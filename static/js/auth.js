function showIn(evt) {
    let none = document.querySelectorAll(".show")
    i = 0
    l = none.length

    for (i; i < l; i++) {
        none[i].style.display = "none";
    }

    let block = document.getElementById(evt.currentTarget.myParam)
    block.style.display = "block"
}

let el1 = document.getElementById("in")
el1.addEventListener("click", showIn, false)
el1.myParam = "login"

let el2 = document.getElementById("up")
el2.addEventListener("click", showIn, false)
el2.myParam = "registration"

let el3 = document.getElementById("formFooter")
el3.addEventListener("click", showIn, false)
el3.myParam = "forget"

/////////////////////////////////////////////////


const form1 = document.getElementById('registration')
form1.onsubmit = async (e) => {
    e.preventDefault()
    await fetch("/registration", {
      method: "POST",
      body: new FormData(form1)
    }).then(function(response) {
        return response.text().then(function(text) {
          if (text == "1") {
            alert("Check Email box and enter login and password")
          } else {
            alert(text)
          }
        })
    })
}

const form2 = document.getElementById('login')
form2.onsubmit = async (e) => {
    e.preventDefault()
    await fetch("/login", {
      method: "POST",
      body: new FormData(form2)
    }).then(function(response) {
        return response.text().then(function(text) {
          if (text == "1") {
            window.location.href = "/adm/allPhones"
          } else {
            alert(text)
          }
        })
    })
}

const form3 = document.getElementById('forget')
form3.onsubmit = async (e) => {
  e.preventDefault()
  await fetch("/forget", {
    method: "POST",
    body: new FormData(form3)
  }).then(function(response) {
      return response.text().then(function(text) {
        if (text == "1") {
          alert("Check Email box and enter login and password")
        } else {
          alert(text)
        }
      })
  })
}




