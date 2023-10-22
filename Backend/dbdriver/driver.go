package dbdriver

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	str "strings"

	"github.com/creasty/defaults"
	"github.com/google/go-querystring/query"
)

// @info All Read-only Operations should use a pointer to the client
// @info All other operations should use copies of this client
type CouchDBClient struct {
	Features    []string `json:"features"`
	Version     string   `json:"version"`
	Sha         string   `json:"git_sha"`
	UUID        string   `json:"uuid"`
	ServerURL   *url.URL
	DatabaseURL *url.URL
	Client      *http.Client
	Vendor      struct {
		Name string `json:"name"`
	} `json:"vendor"`
}

type sizes struct {
	File     uint32 `json:"file"`     // The size of the database file on disk in bytes. Views indexes are not included in the calculation.
	External int32  `json:"external"` // The uncompressed size of database contents in bytes.
	Active   int32  `json:"active"`   // The size of live data inside the database, in bytes.
}

// @info https://docs.couchdb.org/en/stable/api/database/common.html
type DatabaseInfo struct {
	Name      string `json:"db_name"`    // The name of the database.
	PurgeSeq  string `json:"purge_seq"`  // An opaque string that describes the purge state of the database. Do not rely on this string for counting the number of purge operations.
	UpdateSeq string `json:"update_seq"` // An opaque string that describes the state of the database. Do not rely on this string for counting the number of updates.
	Sizes     sizes  `json:"sizes"`
	Props     struct {
		Partitioned bool `json:"partitioned"` // If present and true, this indicates that the database is partitioned.
	} `json:"props"`
	DocDeletedCount   uint32 `json:"doc_del_count"`       // Number of deleted documents
	DocCount          uint32 `json:"doc_count"`           // A count of the documents in the specified database.
	DiskFormatVersion uint16 `json:"disk_format_version"` // The version of the physical format used for the data when it is stored on disk.
	CompactRunning    bool   `json:"compact_running"`     // Set to true if the database compaction routine is operating on this database.
	Cluster           struct {
		Shards      uint16 `json:"q"` // Shards. The number of range partitions.
		Replicas    uint16 `json:"n"` // Replicas. The number of copies of every document.
		WriteQuorum uint16 `json:"w"` // Write quorum. The number of copies of a document that need to be written before a successful reply.
		ReadQuorum  uint16 `json:"r"` // Read quorum. The number of consistent copies of a document that need to be read before a successful reply.
	} `json:"cluster"`
	InstanceStartTime string `json:"instance_start_time"`
}

type View struct {
	MapFunction    string `json:"map"`
	ReduceFunction string `json:"reduce,omitempty"`
}

// @info https://docs.couchdb.org/en/3.2.2/json-structure.html#design-document
type DesignDocument struct {
	ID       string          `json:"_id"`
	REV      string          `json:"_rev"`
	Views    map[string]View `json:"views"`
	Language string          `json:"language,omitempty"`
}

// @info https://docs.couchdb.org/en/3.2.2/json-structure.html#design-document-information
type DesignDocumentInfo struct {
	Name      string `json:"name"`
	ViewIndex struct {
		UpdatesPending struct {
			Minimum   int16 `json:"minimum"`
			Preferred int16 `json:"preferred"`
			Total     int16 `json:"total"`
		} `json:"updates_pending"`
		WaitingCommit  bool   `json:"waiting_commit"`
		WaitingClients int32  `json:"waiting_clients"`
		UpdaterRunning bool   `json:"updater_running"`
		UpdateSeq      uint32 `json:"update_seq"`
		Sizes          sizes  `json:"sizes"`
		Signature      string `json:"signature"`
		PurgeSeq       uint32 `json:"purge_seq"`
		Language       string `json:"javascript"`
		CompactRunning bool   `json:"compact_running"`
	} `json:"view_index"`
}

type GenericDocument = map[string]interface{}

// @info https://docs.couchdb.org/en/3.2.2/api/ddoc/views.html

type Row struct {
	ID    string          `json:"id"`
	Key   GenericDocument `json:"key"`
	Value uint32          `json:"value"`
	Doc   GenericDocument `json:"doc,omitempty"`
}
type DesignView struct {
	TotalRows uint64 `json:"total_rows,omitempty"`
	Offset    uint32 `json:"offset,omitempty"`
	Rows      []Row  `json:"rows"`
}

// https://docs.couchdb.org/en/3.2.2/api/ddoc/views.html
// @info Some of the strings should be JSON (interface{}) but the go-querystring lib can't deal with that so user will have to deal with it manually.
// They are user specific, and so it is the user that, upon need, should create a string that can accomodate it
type designViewOptions struct {
	Conflicts          bool     `default:"false" url:"conflicts,omitempty"`
	Descending         bool     `default:"false" url:"descending,omitempty"`
	EndKey             string   `default:"-" url:"endkey,omitempty"` // @info This is a user-specific JSON
	EndKeyDocID        string   `default:"-" url:"endkey_docid,omitempty"`
	Group              bool     `default:"false" url:"group,omitempty"`
	GroupLevel         uint16   `default:"0" url:"group_level,omitempty"`
	IncludeDocs        bool     `default:"false" url:"include_docs,omitempty"`
	Attachments        bool     `default:"false" url:"attachments,omitempty"`
	AttachEncodingInfo bool     `default:"false" url:"att_encoding_info,omitempty"`
	InclusiveEnd       bool     `default:"false" url:"inclusive_end,omitempty"`
	Key                string   `default:"-" url:"key,omitempty"`   // @info This is a user-specific JSON
	Keys               []string `default:"[]" url:"keys,omitempty"` // @info This is a combination of user-specific JSON files. Multiple entries should exist.
	Limit              uint64   `default:"0" url:"limit,omitempty"`
	Reduce             bool     `default:"false" url:"reduce,omitempty"`
	Skip               uint64   `default:"0" url:"skip,omitempty"`
	Sorted             bool     `default:"true" url:"sorted"`
	Stable             bool     `default:"false" url:"stable,omitempty"`
	StartKey           string   `default:"-" url:"startkey,omitempty"` // @info This is a user-specific JSON
	StartKeyDocID      string   `default:"-" url:"startkey_docid,omitempty"`
	Update             string   `default:"-" url:"update,omitempty"`
	UpdateSeq          bool     `default:"false" url:"update,omitempty"`
	// Stale              string // @info This is deprecated, commented for simplicity. Uncomment and use at your own risk! - Kev 22
}

type PutResponseData struct {
	ID  string `json:"id"`
	OK  bool   `json:"ok"`
	REV string `json:"rev"`
}

// @info https://docs.couchdb.org/en/3.2.2/api/database/find.html
type FindOptions struct {
	Selector       map[string]interface{} `json:"selector,omitempty"` // JSON
	Limit          uint64                 `json:"limit,omitempty"`    // Default in CouchDB is 25
	Skip           uint64                 `json:"skip,omitempty"`
	Sort           []string               `json:"sort,omitempty"` // JSON ARRAY
	Fields         []string               `json:"fields,omitempty"`
	UseIndex       []string               `json:"use_index,omitempty"`
	Conflicts      bool                   `json:"conflicts,omitempty"`
	ReadQuorum     uint64                 `json:"r,omitempty"`
	Bookmark       string                 `json:"bookmark,omitempty"`
	Update         bool                   `json:"update"`
	Stable         bool                   `json:"stable,omitempty"`
	ExecutionStats bool                   `json:"execution_stats,omitempty"`
	// Stale      string                 `default:"-" json:"stale,omitempty"` // @info Deprecated, uncomment and use at own risk
}

type FindResponseData struct {
	Docs           []GenericDocument      `json:"docs"`
	Warning        string                 `json:"warning"`
	ExecutionStats map[string]interface{} `json:"execution_stats,omitempty"`
	Bookmark       string                 `json:"bookmark,omitempty"`
}

func FindInDatabase(client *CouchDBClient, opts *FindOptions) (*FindResponseData, error) {
	resp_data := &FindResponseData{}
	if client.DatabaseURL == nil {
		return resp_data, errors.New("Attempted to find in a database but no database was specified (client is not connected)")
	}
	url := client.DatabaseURL.JoinPath("_find")

	jsonBytes, _ := json.Marshal(&opts)

	fmt.Println("BYTES ", string(jsonBytes))

	req, err := http.NewRequest(http.MethodPost, url.String(), bytes.NewReader(jsonBytes))
	// req, err := http.NewRequest(http.MethodPost, url.String(), bytes.NewReader(buff))
	if err != nil {
		return resp_data, err
	}

	req.Header.Set("Content-Type", "application/json")
	// fmt.Printf("STIUFFFF %v\n", req.Body)

	resp, err := client.Client.Do(req)
	// bodyBytes, err := ioutil.ReadAll(resp.Body)
	// fmt.Println("HMMMMM", string(bodyBytes))
	if err != nil || resp.StatusCode != 200 {
		return resp_data, err
	}
	err = json.NewDecoder(resp.Body).Decode(&resp_data)
	return resp_data, err
}

// url := "http://localhost:5984/{DB_NAME}"
func CreateNewDatabase(client *CouchDBClient, dbname string) string {
	url := client.ServerURL.JoinPath(dbname)
	req, err := http.NewRequest(http.MethodPut, url.String(), bytes.NewBuffer(make([]byte, 0)))
	if err != nil {
		return err.Error()
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Client.Do(req)
	if err != nil || (resp.StatusCode != 201 && resp.StatusCode != 200) {
		return err.Error()
	}
	var ok bool
	fmt.Println(resp.Body)
	err = json.NewDecoder(resp.Body).Decode(&ok)
	return fmt.Sprintf("%t", ok)
}

// url := "http://localhost:5984/{DB_NAME}/{DOC_ID}"
// @opt Reduce String allocations!!! -Kiw 22
// @info Creating and Modifying is exactly the same thing in CouchDB.
func CreateOrModifyDocument(client *CouchDBClient, doc *GenericDocument, id string) (*PutResponseData, error) {
	resp_data := &PutResponseData{}
	if client.DatabaseURL == nil {
		return resp_data, errors.New("Attempted to create a design document from an unspecified database (client is not connected)")
	}
	var err error
	if id == "" {
		id, err = GetUUIDFromCouchDB(client)
		if err != nil {
			return resp_data, err
		}
	}
	url := client.DatabaseURL.JoinPath(id)
	(*doc)["_id"] = id
	data, err := json.Marshal(*doc)
	if err != nil {
		return resp_data, err
	}
	buff := bytes.NewBuffer(data)
	req, err := http.NewRequest(http.MethodPut, url.String(), buff)
	if err != nil {
		return resp_data, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Client.Do(req)
	if err != nil || (resp.StatusCode != 201 && resp.StatusCode != 200) {
		return resp_data, err
	}

	err = json.NewDecoder(resp.Body).Decode(&resp_data)
	return resp_data, err // @error If there's an error it is automatically returned, client must error handle
}

// url := "http://localhost:5984/_uuids"
// @info Alternatively a user may want to resort to Golang's UUID library
func GetUUIDFromCouchDB(client *CouchDBClient) (string, error) {
	type UUIDs struct {
		UUID []string `json:"uuids"`
	}
	id := &UUIDs{}
	url := client.ServerURL.JoinPath("_uuids")
	r, err := http.Get(url.String())
	if err != nil || r.StatusCode != 200 {
		return "", err
	}
	err = json.NewDecoder(r.Body).Decode(&id)
	return id.UUID[0], err // @error If there's an error it is automatically returned, client must error handle
}

func CreateDesignViewOptions() *designViewOptions {
	opts := &designViewOptions{}
	defaults.Set(opts)
	return opts
}

func designViewOptionsToQueryString(opts *designViewOptions) (string, error) {

	if opts.Update != "" && opts.Update != "true" && opts.Update != "false" && opts.Update != "lazy" {
		fmt.Fprintf(os.Stderr, "Warning: Using '%s' as value for DesignViewOption 'Update'. Expected 'true', 'false', or 'lazy'. Please, refer to the documentation at https://docs.couchdb.org/en/3.2.2/api/ddoc/views.html for more information.\n", opts.Update)
	}

	vals, err := query.Values(opts)
	if err != nil {
		return "", err
	}
	if opts.Sorted {
		vals.Del("sorted") // Remove unnecessary param from query, by default it is already true
	}
	if len(opts.Keys) > 0 {
		str := opts.Keys[0]
		for i := 1; i < len(opts.Keys); i++ {
			str += "," + opts.Keys[i]
		}
		vals.Set("keys", str) // Substitute params from query because they already exist
	}
	return vals.Encode(), err
}

// url := "http://localhost:5984/{DB_NAME}/_design/{DESIGN_DOC_NAME}/_info"
// @opt Reduce String allocations!!! -Kiw 22
func GetDesignView(client *CouchDBClient, designDoc string, viewName string, opts *designViewOptions) (*DesignView, error) {
	designView := &DesignView{}
	if client.DatabaseURL == nil {
		return designView, errors.New("Attempted to get a design view from an unspecified database (client is not connected)")
	}
	if !str.HasPrefix(designDoc, "_design/") {
		designDoc = "_design/" + designDoc
	}

	if !str.HasPrefix(viewName, "_view/") {
		viewName = "_view/" + viewName
	}
	url := client.DatabaseURL.JoinPath(designDoc).JoinPath(viewName)

	if opts != nil {
		queryString, err := designViewOptionsToQueryString(opts)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Warning: An error ocurred converting DesignViewOptions to a Query String.")
		}
		url.RawQuery = queryString
	}
	// fmt.Println(url.String())
	r, err := http.Get(url.String())
	if err != nil || r.StatusCode != 200 {
		return designView, err
	}
	err = json.NewDecoder(r.Body).Decode(&designView)
	return designView, err // @error If there's an error it is automatically returned, client must error handle
}

// url := "http://localhost:5984/{DB_NAME}/{DOC_NAME}"
// @opt Reduce String allocations!!! -Kiw 22
func GetDocument(client *CouchDBClient, id string) (GenericDocument, error) {
	var doc = GenericDocument{}
	if client.DatabaseURL == nil {
		return doc, errors.New("Attempted to get a document from an unspecified database (client is not connected)")
	}
	url := client.DatabaseURL.JoinPath(id)
	r, err := http.Get(url.String())
	if err != nil || r.StatusCode != 200 {
		return doc, err
	}
	err = json.NewDecoder(r.Body).Decode(&doc)
	return doc, err // @error If there's an error it is automatically returned, client must error handle
}

// url := "http://localhost:5984/{DB_NAME}/_design/{DESIGN_DOC_NAME}/_info"
// @opt Reduce String allocations!!! -Kiw 22
func GetDesignDocumentInfo(client *CouchDBClient, docname string) (*DesignDocumentInfo, error) {
	docInfo := &DesignDocumentInfo{}
	if client.DatabaseURL == nil {
		return docInfo, errors.New("Attempted to get design document information from an unspecified database (client is not connected)")
	}
	if !str.HasPrefix(docname, "_design/") {
		docname = "_design/" + docname
	}
	docname = docname + "/_info"

	url := client.DatabaseURL.JoinPath(docname)

	r, err := http.Get(url.String())
	if err != nil || r.StatusCode != 200 {
		return docInfo, err
	}
	err = json.NewDecoder(r.Body).Decode(&docInfo)
	return docInfo, err // @error If there's an error it is automatically returned, client must error handle
}

// url := "http://localhost:5984/{DB_NAME}/_design/{DESIGN_DOC_NAME}"
// @opt Reduce String allocations!!! -Kiw 22
func GetDesignDocument(client *CouchDBClient, docname string) (*DesignDocument, error) {
	doc := &DesignDocument{}
	if client.DatabaseURL == nil {
		return doc, errors.New("Attempted to get a design document from an unspecified database (client is not connected)")
	}
	if !str.HasPrefix(docname, "_design/") {
		docname = "_design/" + docname
	}

	url := client.DatabaseURL.JoinPath(docname)

	r, err := http.Get(url.String())
	if err != nil || r.StatusCode != 200 {
		return doc, err
	}
	err = json.NewDecoder(r.Body).Decode(&doc)
	return doc, err // @error If there's an error it is automatically returned, client must error handle
}

// url := "http://localhost:5984/{DB_NAME}"
func getDBInfo(client *CouchDBClient) (*DatabaseInfo, error) {
	db := &DatabaseInfo{}
	r, err := http.Get(client.DatabaseURL.String())
	if err != nil || r.StatusCode != 200 {
		return db, err
	}
	err = json.NewDecoder(r.Body).Decode(&db)
	return db, err // @error If there's an error it is automatically returned, client must error handle
}

// url := "http://localhost:5984/_all_dbs"
func GetAllDBs(client *CouchDBClient) ([]string, error) {
	var dbs []string
	url := client.ServerURL.JoinPath(os.Getenv("ALL_DBS_URL"))
	r, err := http.Get(url.String())
	if err != nil || r.StatusCode != 200 {
		return dbs, err
	}
	err = json.NewDecoder(r.Body).Decode(&dbs)
	return dbs, err // @error If there's an error it is automatically returned, client must error handle
}

func ConnectToDB(client *CouchDBClient, dbname string) (*DatabaseInfo, error) {
	client.DatabaseURL = client.ServerURL.JoinPath(dbname)
	return getDBInfo(client)
}

// url := "http://localhost:5984/"
func CreateClient() (*CouchDBClient, error) {
	url_str := fmt.Sprintf("%s://%s:%s@%s", os.Getenv("COUCHDB_SCH"), os.Getenv("COUCHDB_USR"), os.Getenv("COUCHDB_PWD"), os.Getenv("COUCHDB_URL"))
	new_client := CouchDBClient{}
	var err error
	new_client.ServerURL, err = url.Parse(url_str)
	new_client.Client = &http.Client{}
	// new_client.URL.User = url.UserPassword(os.Getenv("COUCHDB_USR"), os.Getenv("COUCHDB_PWD")) // @info Should not be used unless its for compatibility reasons?
	if err != nil {
		return &new_client, err
	}
	r, err := http.Get(url_str)
	if err != nil {
		return &new_client, err
	}
	err = json.NewDecoder(r.Body).Decode(&new_client)
	if err != nil {
		return &new_client, err
	}
	return &new_client, err
}
