var link = document.createElement( 'link' );
link.rel = 'stylesheet';
link.type = 'text/css';
link.href = window.location.search.match( /print-pdf/gi ) ? 'css/print/pdf.css' : 'css/print/paper.css';
document.getElementsByTagName( 'head' )[0].appendChild( link );
// Websocket
var jarvis = document.getElementById("jarvis");

ws = new WebSocket("wss://localhost:8443/ws");
ws.onopen = function(evt) {
  console.log("OPEN");
};
ws.onclose = function(evt) {
  console.log("CLOSE");
  ws = null;
};
ws.onmessage = function(evt) {
  var msg = new SpeechSynthesisUtterance(evt.data);
  msg.lang = 'fr-FR';
  window.speechSynthesis.speak(msg);
  console.log("RESPONSE: " + evt.data);
  if (evt.data !== "") {
    console.log("Adding bottom");
      document.getElementById("jarvisimg").className = "bottom";
      document.getElementById("jarvisReply").textContent= evt.data;

  } else {
    console.log("Removing bottom");
      document.getElementById("jarvisimg").classList.remove('bottom');

  }
};
ws.onerror = function(evt) {
  console.log("ERROR: " + evt.data);
};



if (annyang) {
  annyang.debug(true);
  // Let's define our first command. First the text we expect, and then the function it should call
  annyang.setLanguage('fr-FR');
  annyang.addCallback('result', function(phrases) {
    ws.send(JSON.stringify(phrases));
    var stop = false;
    for (s of phrases) {
      console.log(s);
      str = s.toLowerCase();
      switch (str) {
        case (str.match(/métrique/) || {}).input:
          console.log("Matched métrique'");
        case (str.match(/monitoring/) || {}).input:
          console.log("Matched the 'test' substring");        
          document.getElementById('metrics').style.visibility = 'visible';
          stop = true;
          break;
        case (str.match(/amélioration/) || {}).input:
          document.getElementById('amcont').style.visibility = 'visible';
          document.getElementById("b_amcont").classList.toggle('blink');
          setTimeout(function () {
            document.getElementById("b_amcont").classList.remove('blink');
          }, 1500);
          stop = true;
          break;
        case (str.match(/continue/) || {}).input:
          document.getElementById('amcont').style.visibility = 'visible';
          document.getElementById("b_amcont").classList.toggle('blink');
          setTimeout(function () {
            document.getElementById("b_amcont").classList.remove('blink');
          }, 1500);
          stop = true;
          break;
        case (str.match(/Scrum/) || {}).input:
          document.getElementById('scrum').style.visibility = 'visible';
          document.getElementById("b_scrum").classList.toggle('blink');
          setTimeout(function () {
            document.getElementById("b_scrum").classList.remove('blink');
          }, 1500);
          stop = true;
          break;
        case (str.match(/string/) || {}).input:
          document.getElementById('sprint').style.visibility = 'visible';
          document.getElementById("b_sprint").classList.toggle('blink');
          setTimeout(function () {
            document.getElementById("b_sprint").classList.remove('blink');
          }, 1500);
          stop = true;
          break;
        case (str.match(/sprint/) || {}).input:
          document.getElementById('sprint').style.visibility = 'visible';
          document.getElementById("b_sprint").classList.toggle('blink');
          setTimeout(function () {
            document.getElementById("b_sprint").classList.remove('blink');
          }, 1500);
          stop = true;
          break;
        case (str.match(/cycle court/) || {}).input:
          document.getElementById('sprint').style.visibility = 'visible';
          document.getElementById("b_sprint").classList.toggle('blink');
          setTimeout(function () {
            document.getElementById("b_sprint").classList.remove('blink');
          }, 1500);
          stop = true;
          break;
        case (str.match(/rapidement/) || {}).input:
          document.getElementById('sprint').style.visibility = 'visible';
          document.getElementById("b_sprint").classList.toggle('blink');
          setTimeout(function () {
            document.getElementById("b_sprint").classList.remove('blink');
          }, 1500);
          stop = true;
          break;
        case (str.match(/rapidité/) || {}).input:
          document.getElementById('sprint').style.visibility = 'visible';
          document.getElementById("b_sprint").classList.toggle('blink');
          setTimeout(function () {
            document.getElementById("b_sprint").classList.remove('blink');
          }, 1500);
          stop = true;
          break;
        case (str.match(/communication/) || {}).input:
          document.getElementById('communication').style.visibility = 'visible';
          console.log("la communication");
          //document.getElementById('b_communication').className = 'blink';
          document.getElementById("b_communication").classList.toggle('blink');
          setTimeout(function () {
            document.getElementById("b_communication").classList.remove('blink');
          }, 1500);
          stop = true;
          break;
        case (str.match(/adaptation/) || {}).input:
          document.getElementById('adaptation').style.visibility = 'visible';
          document.getElementById("b_adaptation").classList.toggle('blink');
          setTimeout(function () {
            document.getElementById("b_adaptation").classList.remove('blink');
          }, 1500);
          stop = true;
          break;
        case (str.match(/adaptabilité/) || {}).input:
          document.getElementById('adaptation').style.visibility = 'visible';
          document.getElementById("b_adaptation").classList.toggle('blink');
          setTimeout(function () {
            document.getElementById("b_adaptation").classList.remove('blink');
          }, 1500);
          stop = true;
          break;
        case (str.match(/produit/) || {}).input:
          document.getElementById('produit').style.visibility = 'visible';
          document.getElementById("b_produit").classList.toggle('blink');
          setTimeout(function () {
            document.getElementById("b_produit").classList.remove('blink');
          }, 1500);
          stop = true;
          break;
        case (str.match(/projet/) || {}).input:
          document.getElementById('produit').style.visibility = 'visible';
          document.getElementById("b_produit").classList.toggle('blink');
          setTimeout(function () {
            document.getElementById("b_produit").classList.remove('blink');
          }, 1500);
          stop = true;
          break;
        case (str.match(/méthode/) || {}).input:
          document.getElementById('scrum').style.visibility = 'visible';
          document.getElementById("b_scrum").classList.toggle('blink');
          setTimeout(function () {
            document.getElementById("b_scrum").classList.remove('blink');
          }, 1500);
          stop = true;
          break;
        case (str.match(/collabore²/) || {}).input:
          document.getElementById('equipe').style.visibility = 'visible';
          document.getElementById("b_equipe").classList.toggle('blink');
          setTimeout(function () {
            document.getElementById("b_equipe").classList.remove('blink');
          }, 1500);
          stop = true;
          break;
        case (str.match(/collaboration/) || {}).input:
          document.getElementById('equipe').style.visibility = 'visible';
          document.getElementById("b_equipe").classList.toggle('blink');
          setTimeout(function () {
            document.getElementById("b_equipe").classList.remove('blink');
          }, 1500);
          stop = true;
          break;
        case (str.match(/équipe/) || {}).input:
          document.getElementById('equipe').style.visibility = 'visible';
          document.getElementById("b_equipe").classList.toggle('blink');
          setTimeout(function () {
            document.getElementById("b_equipe").classList.remove('blink');
          }, 1500);
          stop = true;
          break;
        case (str.match(/cloud/) || {}).input:
          document.getElementById('agile').style.visibility = 'visible';
          document.getElementById('devops').style.visibility = 'visible';
          stop = true;
          break;
        case (str.match(/agile/) || {}).input:
          document.getElementById('agile').style.visibility = 'visible';
          document.getElementById('devops').style.visibility = 'visible';
          stop = true;
          break;
        case (str.match(/agilité/) || {}).input:
          document.getElementById('agile').style.visibility = 'visible';
          document.getElementById('devops').style.visibility = 'visible';
          stop = true;
          break;
        case (str.match(/gilles/) || {}).input:
          document.getElementById('agile').style.visibility = 'visible';
          document.getElementById('devops').style.visibility = 'visible';
          stop = true;
          break;
        default:
          console.log("Didn't match");
          stop = true;
          break;
      }
      if (stop)  {
        break;
      }

      //if (s.toLowerCase().indexOf("metrique") >= 0 || s.toLowerCase().indexOf("monitoring")) {
      //    document.getElementById('metrics').style.visibility = 'visible';
      //} else if (s.toLowerCase().indexOf("corps") >= 0 || s.toLowerCase().indexOf("edge") >= 0 || s.toLowerCase().indexOf("Corée")) {
      //    document.getElementById('coreedge').style.visibility = 'visible';
      //} else if (s.toLowerCase().indexOf("amélioration continue") >= 0) {
      //    document.getElementById('amcont').style.visibility = 'visible';
      //}

    }
    //console.log("I think the user said: ", phrases[0]);
    //  console.log("But then again, it could be any of the following: ", phrases);

  });

  // Start listening. You can call this here, or attach this call to an event, button, etc.
  annyang.start();
}

