// $("#host-submit").on("click", (e)=>{
//     e.preventDefault()

//     email=$("#email").val()
//     password=$("#password").val()
//     port1=$("#port1").val()
//     port2=$("#port2").val()
//     image1=$("#image1").val()
//     image2=$("#image2").val()


//     $.post("/host", JSON.stringify({
//         email,
//         password,
//         port1,
//         port2,
//         image1,
//         image2
//     }), (data)=>{
//         if(data) {
//             localStorage.setItem("port1", port1);
//             localStorage.setItem("image1", image1);
//             localStorage.setItem("port2", port2);
//             localStorage.setItem("image2", image2);
//             window.location.replace("/session")
//         }
//     })
// })