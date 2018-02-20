import 'materialize-css/dist/css/materialize.min.css'
import anime from 'animejs'
import '../css/style.css';
import mountain from '../assets/mountain-min.png';
import rmountain from '../assets/rmountain-min.png';
import me from '../assets/me3.png';


(function () {
    console.log(process.env.NODE_ENV);

    //document.getElementById('backgroundOne').style.backgroundImage = "url('" + mountain + "')";
    //document.getElementById('backgroundTwo').style.backgroundImage = "url('" + rmountain + "')";
    $(document).ready(function () {
        $('.parallax').parallax();
    });
})();