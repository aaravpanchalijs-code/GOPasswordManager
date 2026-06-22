const API_URL = "http://localhost:8080";

const token = localStorage.getItem("token");

if (!token) {
  window.location.href = "login.html";
}

const passwordGrid = document.getElementById("passwordGrid");
const addPasswordForm = document.getElementById("addPasswordForm");
const logoutBtn = document.getElementById("logoutBtn");
const toast = document.getElementById("toast");

let currentEditID = "";
let currentShareID = "";

// ============================
// Toast
// ============================

function showToast(message, success = true) {
  toast.innerText = message;

  toast.style.borderLeftColor = success ? "#22C55E" : "#EF4444";

  toast.classList.add("show");

  setTimeout(() => {
    toast.classList.remove("show");
  }, 3000);
}

// ============================
// Logout
// ============================

logoutBtn.addEventListener("click", () => {
  localStorage.removeItem("token");

  window.location.href = "login.html";
});

// ============================
// Load Passwords
// ============================

async function loadPasswords() {
  passwordGrid.innerHTML = "<h3>Loading...</h3>";

  try {
    const response = await fetch(`${API_URL}/vault/get`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });

    let passwords;

    try {
      passwords = await response.json();
    } catch {
      showToast("Invalid server response", false);

      return;
    }

    if (!response.ok) {
      showToast(passwords.message, false);

      return;
    }

    renderPasswords(passwords);
  } catch (err) {
    console.error(err);

    showToast("Server Error", false);
  }
}

// ============================
// Render Cards
// ============================

function renderPasswords(passwords) {
  passwordGrid.innerHTML = "";
  document.getElementById("passwordCount").innerText = passwords.length;
  document.getElementById("sharedCount").innerText = "N/A";

  if (passwords.length === 0) {
    passwordGrid.innerHTML = "<h3>No Passwords Saved Yet</h3>";

    return;
  }

  passwords.forEach((password) => {
    const encodedPassword = encodeURIComponent(password.password);

    const encodedWebsite = encodeURIComponent(password.website);
    const encodedEmail = encodeURIComponent(password.login_email);
    const encodedNotes = encodeURIComponent(password.notes || "");

    const card = document.createElement("div");

    card.className = "password-card";

    card.innerHTML = `

            <h3>${password.website}</h3>

            <p>
                <strong>Email:</strong>
                ${password.login_email}
            </p>

            <p>

                <strong>Password:</strong>

                <span id="pass-${password.id}">
                    ********
                </span>

            </p>

            <p>${password.notes || ""}</p>

            <div class="card-buttons">

                <button
                    class="show"
                    onclick="togglePassword(
                        'pass-${password.id}',
                        '${encodedPassword}'
                    )">

                    Show

                </button>

               <button
    class="edit"
    onclick="openEditModal(
        '${password.id}',
        '${encodedWebsite}',
        '${encodedEmail}',
        '${encodedPassword}',
        '${encodedNotes}'
    )">

    Edit

</button>

                <button
                    class="share"
                    onclick="openShareModal('${password.id}')">

                    Share

                </button>

                <button
                    class="delete"
                    onclick="deletePassword('${password.id}')">

                    Delete

                </button>

            </div>

        `;

    passwordGrid.appendChild(card);
  });
}

// ============================
// Show Password
// ============================

function togglePassword(id, encodedPassword) {
  const span = document.getElementById(id);

  const password = decodeURIComponent(encodedPassword);

  if (span.innerText === "********") {
    span.innerText = password;
  } else {
    span.innerText = "********";
  }
}

// ============================
// Add Password
// ============================

addPasswordForm.addEventListener("submit", async (e) => {
  e.preventDefault();

  const website = document.getElementById("website").value.trim();

  const loginEmail = document.getElementById("loginEmail").value.trim();

  const password = document.getElementById("password").value;

  const notes = document.getElementById("notes").value.trim();

  try {
    const response = await fetch(`${API_URL}/vault/add`, {
      method: "POST",

      headers: {
        "Content-Type": "application/json",

        Authorization: `Bearer ${token}`,
      },

      body: JSON.stringify({
        website,

        login_email: loginEmail,

        password,

        notes,
      }),
    });

    const data = await response.json();

    if (!response.ok) {
      showToast(data.message, false);

      return;
    }

    showToast(data.message);

    addPasswordForm.reset();

    loadPasswords();
  } catch (err) {
    console.error(err);

    showToast("Server Error", false);
  }
});

// Initial Load

loadPasswords();


// =====================================
// DELETE PASSWORD
// =====================================

async function deletePassword(id){

    if(!confirm("Are you sure you want to delete this password?")){
        return;
    }

    try{

        const response = await fetch(
            `${API_URL}/vault/delete?id=${id}`,
            {
                method:"DELETE",
                headers:{
                    Authorization:`Bearer ${token}`
                }
            }
        );

        const data = await response.json();

        if(!response.ok){

            showToast(data.message,false);
            return;

        }

        showToast(data.message);

        loadPasswords();

    }
    catch(err){

        console.error(err);

        showToast("Server Error",false);

    }

}

// =====================================
// EDIT MODAL
// =====================================

const editModal = document.getElementById("editModal");

const editWebsite = document.getElementById("editWebsite");

const editEmail = document.getElementById("editEmail");

const editPassword = document.getElementById("editPassword");

const editNotes = document.getElementById("editNotes");

const updateBtn = document.getElementById("updateBtn");

function openEditModal(
    id,
    website,
    email,
    password,
    notes
){

    currentEditID = id;

    editWebsite.value = decodeURIComponent(website);

    editEmail.value = decodeURIComponent(email);

    editPassword.value = decodeURIComponent(password);

    editNotes.value = decodeURIComponent(notes);

    editModal.classList.add("show");

}

// Close Edit Modal

window.addEventListener("click",(e)=>{

    if(e.target===editModal){

        editModal.classList.remove("show");

    }

});

// =====================================
// UPDATE PASSWORD
// =====================================

updateBtn.addEventListener("click",async()=>{

    try{

        const response = await fetch(

            `${API_URL}/vault/update?id=${currentEditID}`,

            {

                method:"PUT",

                headers:{

                    "Content-Type":"application/json",

                    Authorization:`Bearer ${token}`

                },

                body:JSON.stringify({

                    website:editWebsite.value,

                    login_email:editEmail.value,

                    password:editPassword.value,

                    notes:editNotes.value

                })

            }

        );

        const data = await response.json();

        if(!response.ok){

            showToast(data.message,false);

            return;

        }

        showToast(data.message);

        editModal.classList.remove("show");

        loadPasswords();

    }
    catch(err){

        console.error(err);

        showToast("Server Error",false);

    }

});

// =====================================
// SHARE MODAL
// =====================================

const shareModal = document.getElementById("shareModal");

const shareUsername = document.getElementById("shareUsername");

const shareBtn = document.getElementById("shareBtn");

function openShareModal(id){

    currentShareID = id;

    shareModal.classList.add("show");

}

window.addEventListener("click",(e)=>{

    if(e.target===shareModal){

        shareModal.classList.remove("show");

    }

});

// =====================================
// SHARE PASSWORD
// =====================================

shareBtn.addEventListener("click",async()=>{

    const username = shareUsername.value.trim();

    if(username===""){

        showToast("Enter a username",false);

        return;

    }

    try{

        const response = await fetch(

            `${API_URL}/vault/share`,

            {

                method:"POST",

                headers:{

                    "Content-Type":"application/json",

                    Authorization:`Bearer ${token}`

                },

                body:JSON.stringify({

                    password_id:currentShareID,

                    username

                })

            }

        );

        const data = await response.json();

        if(!response.ok){

            showToast(data.message,false);

            return;

        }

        showToast(data.message);

        shareUsername.value="";

        shareModal.classList.remove("show");

    }
    catch(err){

        console.error(err);

        showToast("Server Error",false);

    }

});


// =======================================
// Generate Strong Password
// =======================================

const generateBtn = document.getElementById("generateBtn");

generateBtn.addEventListener("click", () => {

    const chars =
        "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*()_+";

    let password = "";

    for(let i=0;i<16;i++){

        password += chars.charAt(
            Math.floor(Math.random()*chars.length)
        );

    }

    document.getElementById("password").value = password;

});