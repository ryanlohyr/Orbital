include "base.thrift"
namespace go api

struct Response {
    1: string action
    255: base.BaseResp BaseResp
}

struct ReviewRequest{
    //review 
    1: string Msg 
    //userID
    2: i64 userID
}

struct EditRequest{
    //reviewID
    1: i64 reviewID
    //postID
    2: i64 postID 
    //new review
    3:string Msg
}

struct DeleteRequest{
    //reviewID
    1: i64 reviewID
}

service ReviewService {
    Response sendReview(1: ReviewRequest req)
    Response editReview(1: EditRequest req)
    Response deleteReview(1: DeleteRequest req)
}

