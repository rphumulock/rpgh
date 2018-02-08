/*import 'materialize-css/dist/css/materialize.min.css'
import '../css/style.css';
import Vivus from 'vivus';

(function () {
  window.onload = () => {
    let visited = localStorage.getItem("visited");
    if (visited) {
      window.location.assign("home.html");
    } else {
      localStorage.setItem("visited", "true");
      new Vivus('hi', {
        type: 'delayed',
        duration: 200,
        onReady: function (obj) {
          obj.el.style.opacity = 1;
        }
      }, () => {
        window.location.assign("home.html");
      });
    }
  };
})();

import 'materialize-css/dist/css/materialize.min.css'
import anime from 'animejs'
import '../css/style.css';

(function () {
  console.log(process.env.NODE_ENV);
  window.toggleNav = toggleNav;
  window.enterButton = enterButton;
  window.leaveButton = leaveButton;

  var buttonEl = document.querySelector('.ball');
  buttonEl.addEventListener('mouseenter', enterButton, false);
  buttonEl.addEventListener('mouseleave', leaveButton, false);


  /* let divArr = [];
   for (let i = 0; i < 30; i++) {
       let div = document.createElement("div");
       div.setAttribute('id', "div" + i);
       div.style.width = "28px";
       div.style.height = "28px";
       div.style.backgroundColor = "#FF1461";
       //div.style.opacity = 0;
       div.style.position = "fixed";
       divArr.push(div);
   }
 
   divArr.forEach((div) => {
       console.log(window.height);
       let posX = Math.floor(Math.random() * (window.innerWidth / 10)) + 1;
       let posY = Math.floor(Math.random() * (window.innerHeight / 10)) + 1;
       div.style.left = posX + "%";
       div.style.top = posY + "%";
       document.getElementById("main").appendChild(div);
   });
 
   /*background-size: cover;
     background-position: center;
     background-repeat: no-repeat;
     document.getElementById('main').style.backgroundImage = "url('" + gif + "')";
})();

function toggleNav() {
  if (document.getElementById("mySidenav").style.width === "250px") {
    document.getElementById("mySidenav").style.width = "0";
    document.getElementById("main").style.marginLeft = "0";
  } else {
    document.getElementById("mySidenav").style.width = "250px";
    document.getElementById("main").style.marginLeft = "250px";
  }
}

  background-size: cover;
  background-position: center;
  background-repeat: no-repeat;

function animateButton(scale, duration, elasticity) {
  anime.remove('.ball');
  anime({
    targets: '.ball',
    rotate: {
      value: 360,
      easing: 'easeInOutSine'
    },
    loop: true
  });
}

function enterButton() { animateButton(1.2, 800, 400) };
function leaveButton() { animateButton(1.0, 600, 300) };




    /*var relativeOffset = anime.timeline();

    relativeOffset
        .add({
            targets: '#s0',
            translateX: 250,
            easing: 'easeOutExpo',
        })
        .add({
            targets: '#s1',
            opacity: 1,
            easing: 'easeOutExpo',
            offset: '-=600' // Starts 600ms before the previous animation ends
        })
        .add({
            targets: '#s1',
            opacity: 0,
            easing: 'easeOutExpo',
            offset: '-=800' // Starts 800ms before the previous animation ends
        });*/
