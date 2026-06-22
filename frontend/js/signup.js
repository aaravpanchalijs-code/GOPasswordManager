const signupForm = document.getElementById("signupForm");

const togglePassword = document.getElementById("togglePassword");

const passwordInput = document.getElementById("password");

const toast = document.getElementById("toast");

// ----------------------
// Show / Hide Password
// ----------------------

togglePassword.addEventListener("click", () => {

    if(passwordInput.type === "password"){

        passwordInput.type = "text";

        togglePassword.classList.remove("fa-eye");
        togglePassword.classList.add("fa-eye-slash");

    }
    else{

        passwordInput.type = "password";

        togglePassword.classList.remove("fa-eye-slash");
        togglePassword.classList.add("fa-eye");

    }

});

// ----------------------
// Toast
// ----------------------

function showToast(message, success = true){

    toast.innerText = message;

    toast.style.borderLeftColor = success ? "#22C55E" : "#EF4444";

    toast.classList.add("show");

    setTimeout(()=>{

        toast.classList.remove("show");

    },3000);

}

// ----------------------
// Signup
// ----------------------

signupForm.addEventListener("submit", async function(e){

    e.preventDefault();

    const username = document.getElementById("username").value.trim();

    const email = document.getElementById("email").value.trim();

    const password = passwordInput.value;

    try{

        const response = await fetch("http://localhost:8080/signup",{

            method:"POST",

            headers:{
                "Content-Type":"application/json"
            },

            body:JSON.stringify({

                username,
                email,
                password

            })

        });

        const data = await response.json();

        if(!response.ok){

            showToast(data.message || "Signup Failed",false);

            return;

        }

        showToast("Account Created Successfully!");

        setTimeout(()=>{

            window.location.href = "login.html";

        },1200);

    }

    catch(error){

        console.error(error);

        showToast("Server Error",false);

    }

});