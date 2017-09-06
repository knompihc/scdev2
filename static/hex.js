// $(document).ready(function(){
//     console.log("came here");
//     $('#hexForm').on('submit',function(e){
//         //give hex code request
//         console.log("came here111");
//         $.ajax({
//             type: "GET",
//             url: "/resp",
//             success: function(response){
//                 console.log("response recieved "+response);
//                 $("#response").innerHTML=response;
//             }
//         });
//     });
// });
$(document).ready(function(){
var hexForm = new Vue({
    el: "#hexForm",
    data:{
        message: "just test"
    },
    methods: {
        getResponse: function(e){
            console.log("came here111");
            // $.ajax({
            //     type: "GET",
            //     url: "http://localhost:8004/resp",
            //     success: function(response){
            //         console.log("response recieved "+response);
                    $("#response").innerHTML=response;
            //     }
            // })
        }
    }
})
})