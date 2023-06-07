namespace go api

struct Item{
    // For nested structures, if you want to set the serialization key, use gotag, such as `json: "Id"`
    1: optional i64 id(go.tag = 'json:"id"')
    2: optional string text
}

typedef string JsonDict
struct BizRequest {
    // Corresponding to v_int64 in HTTP query, and the value range is (0, 200)
    //1: Item some(api.body = 'some') // Corresponding first level key = some
	1: string message
  
}

struct BizResponse {
    // This field will be filled in the header returned to the client
    1: string successMessage
   
}

service BizService{
    // Example:   post request
    BizResponse BizMethod1(1: BizRequest req)(
        api.post = '/',
        api.baseurl = 'https://127.0.0.1:8888',
		api.path = 'https://127.0.0.1:8888',
        api.param = 'true',
        api.serializer = 'json'
    )
}

