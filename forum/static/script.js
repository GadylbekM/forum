let navBtn = document.querySelector(".nav_btn")
let sideBar = document.querySelector(".sidebar")
navBtn.addEventListener("click", () => {
  navBtn.classList.toggle("active")
  sideBar.classList.toggle("open")
})






// let reactionBtn = document.querySelectorAll(".reaction-button")
// fn      = function(e) { e.preventDefault() };

// for ( var i = reactionBtn.length; i--; ) {
//   reactionBtn[i].addEventListener('click', fn, false);
// }