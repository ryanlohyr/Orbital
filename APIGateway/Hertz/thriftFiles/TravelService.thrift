include "base.thrift"
namespace go kitex.test.server

struct ClientReq {
    1: required string Msg,
    255: base.Base Base,
}
struct ClientResp {
    1: required string Msg,
    255: base.BaseResp BaseResp,
}
service TravelService {
    ClientReq SendClientData(1: ClientReq req),
}




