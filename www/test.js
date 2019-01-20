const Http = new XMLHttpRequest();
const url = '/posts';
Http.open("GET", url);
Http.send();
Http.onreadystatechange = (e) => {
    console.log(Http.responseText)
}