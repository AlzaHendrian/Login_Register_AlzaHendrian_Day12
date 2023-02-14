// get data arrow function
const getData = () => {
 // data variable
 let name = document.querySelector("#name").value;
 let email = document.querySelector("#email").value;
 let phone = document.querySelector("#phone").value;
 let subject = document.querySelector("#subjects").value;
 let message = document.querySelector("#message").value;

 // alert validate
 // const notif = document.querySelector("#alert").style.display = "block"

 // validate data null
 //   if(name == "") {
 //   //  notif.innerHTML = "nama tidak boleh kosong"
 //   alert('data tidak boleh kosong')
 // }
 // else if(!email) {
 //     alert('data tidak boleh kosong')
 //     // notif.innerHTML = "email tidak boleh kosong"
 //   }
 //   else if(!phone) {
 //     alert('data tidak boleh kosong')
 //     // notif.innerHTML = "nomor telepon tidak boleh kosong"
 //   }
 //   else if(!subject) {
 //     alert('data tidak boleh kosong')
 //     // notif.innerHTML = "subject tidak boleh kosong"
 //   }
 //   else if(!message) {
 //     alert('data tidak boleh kosong')
 //     // notif.innerHTML ="pesan tidak boleh kosong"
 //   }

 //   // to mail
 //   else {
 //     alert('hello')
 //   //   const toEmail = 'tommymh21@gmail.com'
 //   //  let bodymsg = `Halo nama saya ${name}, saya ingin ${description} tolong hubungin saya ${phone}`
 //   //  window.open(`mailto:${toEmail}?subject=${subject}&body=${bodymsg}`)
 //   }

   let data = {
     name,
     email,
     phone,
     subject,
     message
   }
   console.log(data);
}