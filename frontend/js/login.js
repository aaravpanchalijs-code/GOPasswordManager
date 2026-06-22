const loginForm = document.getElementById("loginForm");
const togglePassword = document.getElementById("togglePassword");
const passwordInput = document.getElementById("password");
const toast = document.getElementById("toast");

// ---------------------------
// Show / Hide Password
// ---------------------------

togglePassword.addEventListener("click", () => {

    if(passwordInput.type === "password"){

        passwordInput.type = "text";

        togglePassword.classList.remove("fa-eye");
        togglePassword.classList.add("fa-eye-slash");

    }else{

        passwordInput.type = "password";

        togglePassword.classList.remove("fa-eye-slash");
        togglePassword.classList.add("fa-eye");

    }

});

// ---------------------------
// Toast Notification
// ---------------------------

function showToast(message, success = true){

    toast.textContent = message;

    toast.style.borderLeftColor = success ? "#22C55E" : "#EF4444";

    toast.classList.add("show");

    setTimeout(() => {

        toast.classList.remove("show");

    },3000);

}

// ---------------------------
// Login
// ---------------------------

loginForm.addEventListener("submit", async function(e){

    e.preventDefault();

    const email = document.getElementById("email").value.trim();

    const password = passwordInput.value;

    try{

        const response = await fetch("http://localhost:8080/login",{

            method:"POST",

            headers:{
                "Content-Type":"application/json"
            },

            body:JSON.stringify({

                email,
                password

            })

        });

        const data = await response.json();

        if(!response.ok){

            showToast(data.message || "Login Failed",false);

            return;

        }

        localStorage.setItem("token",data.token);

        showToast("Login Successful!");

        setTimeout(()=>{

            window.location.href="dashboard.html";

        },1000);

    }

    catch(error){

        console.error(error);

        showToast("Server Error",false);

    }

});