include "base.thrift"
namespace go kitex.test.server

struct ClientReq {
    1: required string Msg,
    255: base.Base Base,
}

struct GetClientReq {
    1: required i32 userID,
    255: base.Base Base,
}

struct ClientResp {
    1: required string Msg,
    255: base.BaseResp BaseResp,
}

struct TravelDestResp {
    1: required list<string> Destinations,
    255: base.BaseResp BaseResp,
}

struct RetrieveClientResp{
    1: required string Name,
    2: required i64 id,
    3: required list<string> VisitedCountries,
    255: base.BaseResp BaseResp,

}

service TravelService {
    ClientResp SendClientData(1: ClientReq req),
    RetrieveClientResp RetrieveClientData(1: GetClientReq req),
    TravelDestResp GetAllTravelDestinations(1: GetClientReq req),
}




