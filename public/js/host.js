$("#host-submit").on("click", (e)=>{
    e.preventDefault()
    $.post("/host", JSON.stringify({
        email:$("#email").val(),
        password:$("#password").val(),
        port1:$("#port1").val(),
        port2:$("#port2").val(),
        image1:$("#image1").val(),
        image2:$("#image2").val()
    }), (data)=>{
        console.log(data);
    })
})