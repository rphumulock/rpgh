import 'materialize-css/dist/css/materialize.min.css'
import anime from 'animejs'
import '../css/style.css';
import mountain from '../assets/mountain-min.png';
import desert from '../assets/desert-min.png';
import me from '../assets/me.png';

(function () {
    console.log(process.env.NODE_ENV);
    window.toggleNav = toggleNav;
    window.enterButton = enterButton;
    window.leaveButton = leaveButton;
    window.showBackgroundSlides = showBackgroundSlides;
    window.showProfileSlides = showProfileSlides;

    document.getElementById('backgroundOne').style.backgroundImage = "url('" + mountain + "')";
    document.getElementById('backgroundTwo').style.backgroundImage = "url('" + desert + "')";
    document.getElementById('profile-pic').setAttribute('src', me);

    var els = document.querySelectorAll('#nodeList .el');

    var nodeList = anime({
        targets: els,
        translateX: 50,
        loop: true,
        direction: 'alternate',
        autoplay: true
    });


    var functionBasedDuration = anime({
        targets: '#functionBasedDuration .el',
        translateX: 50,
        direction: 'alternate',
        loop: true,
        duration: function (el, i, l) {
            return 400 + (i * 400);
        }
    });

    var customBezier = anime({
        targets: '#customBezier .el',
        translateX: 50,
        direction: 'alternate',
        easing: [.91, -0.54, 0.29, -0.56],
        rotate: 180,
        loop: true
    });

    var combinedFunctionBasedProp = anime({
        targets: '#combinedFunctionBasedProp .el',
        translateX: 50,
        rotate: 360,
        elasticity: 300,
        duration: function (target) {
            // Duration based on every div 'data-duration' attribute
            return target.getAttribute('data-duration');
        },
        delay: function (target, index) {
            // 100ms delay multiplied by every div index, in ascending order
            return index * 100;
        },
        elasticity: function (target, index, totalTargets) {
            // Elasticity multiplied by every div index, in descending order
            return 200 + ((totalTargets - index) * 200);
        },
        direction: 'alternate',
        loop: true
    });

    window.enterToggle = enterToggle;
    window.leaveToggle = leaveToggle;

    var buttonEl = document.querySelector('.togglenav');
    buttonEl.addEventListener('mouseenter', enterButton, false);
    buttonEl.addEventListener('mouseleave', leaveButton, false);
    buttonEl.addEventListener('click', toggleNav, false);

    var buttonEl = document.querySelector('.toggle');
    buttonEl.addEventListener('mouseenter', enterToggle, false);
    buttonEl.addEventListener('mouseleave', leaveToggle, false);

    window.slideIndex = 0;
    showBackgroundSlides();
})();

function enterToggle() {
    g(true);
}

function leaveToggle() {
    g(false);
}

function g(loop) {
    anime.remove('#keyframes .el');
    var keyframes = anime({
        targets: '#keyframes .el',
        translateX: [
            { value: 50, duration: 1000, delay: 500, elasticity: 0 },
            { value: 0, duration: 1000, delay: 500, elasticity: 0 }
        ],
        translateY: [
            { value: 15, duration: 500, delay: 1000, elasticity: 100 },
            { value: 0, duration: 500, delay: 1000, elasticity: 100 }
        ],
        scaleX: [
            { value: 4, duration: 100, delay: 500, easing: 'easeOutExpo' },
            { value: 1, duration: 900, elasticity: 300 },
            { value: 4, duration: 100, delay: 500, easing: 'easeOutExpo' },
            { value: 1, duration: 900, elasticity: 300 }
        ],
        scaleY: [
            { value: [1.75, 1], duration: 500 },
            { value: 2, duration: 50, delay: 1000, easing: 'easeOutExpo' },
            { value: 1, duration: 450 },
            { value: 1.75, duration: 50, delay: 1000, easing: 'easeOutExpo' },
            { value: 1, duration: 450 }
        ],
        loop: loop
    });
}

function toggleNav() {
    if (document.getElementById("mySidenav").style.width === "250px") {
        document.getElementById("mySidenav").style.width = "0";
        document.getElementById("main").style.marginLeft = "0";
        document.getElementById("details").style.display = "block";
    } else {
        document.getElementById("mySidenav").style.width = "250px";
        document.getElementById("main").style.marginLeft = "250px";
        document.getElementById("details").style.display = "none";
    }
}

function showBackgroundSlides() {
    var i;
    var slides = document.getElementsByClassName("myBackgroundSlides");
    for (i = 0; i < slides.length; i++) { slides[i].style.visibility = "hidden"; }
    slideIndex++;
    if (slideIndex > slides.length) { slideIndex = 1 }
    slides[slideIndex - 1].style.visibility = "visible";
    setTimeout(showBackgroundSlides, 12000);
}

function showProfileSlides() {
    var i;
    var slides = document.getElementsByClassName("myProfileSlides");
    for (i = 0; i < slides.length; i++) { slides[i].style.visibility = "hidden"; }
    slideIndex++;
    if (slideIndex > slides.length) { slideIndex = 1 }
    slides[slideIndex - 1].style.visibility = "visible";
    setTimeout(showProfileSlides, 8000);
}

function animateButton(scale, value, loop) {
    anime.remove('.togglenav');
    anime({
        targets: '.togglenav',
        rotate: {
            value: value,
            easing: 'easeInOutSine'
        },
        direction: 'alternate',
        scale: scale,
        loop: loop
    });
};

function enterButton() {
    animateButton(1.5, 360, true);
};
function leaveButton() {
    animateButton(1, 0, false);
};









    /*
      //document.getElementById('background').style.backgroundImage = "url('" + gif + "')";

    var timelineParameters = anime.timeline({
        loop: true,
        easing: 'easeInOutSine'
    });

    timelineParameters
        .add({
            targets: '#background',
            scale: 1.2,
            duration: 30000
        })
        .add({
            targets: '#background',
            scale: 1,
            duration: 30000
        })
    
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
   
     background-size: cover;
       background-position: center;
       background-repeat: no-repeat;
       document.getElementById('main').style.backgroundImage = "url('" + gif + "')";
    
    
    
    var relativeOffset = anime.timeline();

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
