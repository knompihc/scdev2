$(document).ready(function(){
    console.log("came here");
    $('#hexForm').on('submit',function(e){
        //give hex code request
        console.log("came here111");
        $.ajax({
            type: "GET",
            url: "/resp",
            success: function(response){
                console.log("response recieved "+response);
                $("#response").innerHTML=response;
            }
        });
    });
});