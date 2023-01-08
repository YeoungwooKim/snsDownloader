
// document.getElementById("inputUrl").placeholder = 'https://twitter.com/i/status/1273993946406907904'
document.getElementById("inputUrl").value = "https://twitter.com/NASA/status/1606686673915584512"
//"https://www.youtube.com/watch?v=C0DPdy98e4c&ab_channel=SimonYapp"
//"https://twitter.com/MUPLY_PLAYLIST/status/1273993946406907904"


var response = null;

document.getElementById('bestQualityCheckBox').addEventListener('click', function (event) {
    downloadBtn = document.getElementById("downloadButton")
    toggleBtn = document.getElementById("toggleButton")
    searchBtn = document.getElementById("searchButton")
    if (qualityCheckBox.checked) {
        console.log("hello checked..")
        searchBtn.value = "Download"
        swapStatus(downloadBtn, toggleBtn)
    } else {
        console.log("hello not checked...")
        searchBtn.value = "Search"
        swapStatus(downloadBtn, toggleBtn)
    }
    searchBtn.innerHTML = searchBtn.value
    console.log(searchBtn)
})

function swapStatus(downloadBtn, toggleBtn) {
    document.getElementById("progressBar").style.display = 'none'
    if (downloadBtn == null || toggleBtn == null) {
        console.log("both is null")
        return
    }
    if (toggleBtn.style.display == "") {
        toggleBtn.style.display = 'none';
        downloadBtn.style.display = 'none';
        return
    }
    toggleBtn.style.display = '';
    downloadBtn.style.display = '';
}

function btnDownload() {
    var obj = new Object();
    obj.platform = document.getElementById("inputGroupSelect").value;
    obj.uri = document.getElementById("inputUrl").value;
    document.getElementById("progressBar").style.display = "";

    videoOption = document.querySelector("input[name=videoOptions]:checked")
    audioOption = document.querySelector("input[name=audioOptions]:checked")

    if (videoOption != null) {
        obj.videoId = videoOption.value
    }
    if (audioOption != null) {
        obj.audioId = audioOption.value
    }

    $.ajax({
        method: 'post',
        async: false, // 동기 요청으로 변경
        headers: {
            'Content-Type': 'application/json;charset=utf-8'
        },

    })

    handleFileDownload('/api/v1/media', obj)
}

async function handleFileDownload(url, requestBody) {
    const response = await fetch(url, {
        method: 'post',
        headers: {
            'Content-Type': 'application/json;charset=utf-8'
        },
        body: JSON.stringify(requestBody)
    });
    const file = await response.blob();
    fileName = response.headers.get('content-disposition') + ""
    fileName = fileName.replace("attachment; filename=", "")
    fileName = fileName.replaceAll(`"`, '')
    const downloadUrl = window.URL.createObjectURL(file); // 해당 file을 가리키는 url 생성

    const anchorElement = document.createElement('a');
    document.body.appendChild(anchorElement);
    anchorElement.download = fileName; // a tag에 download 속성을 줘서 클릭할 때 다운로드가 일어날 수 있도록 하기
    anchorElement.href = downloadUrl; // href에 url 달아주기

    anchorElement.click(); // 코드 상으로 클릭을 해줘서 다운로드를 트리거

    document.body.removeChild(anchorElement); // cleanup - 쓰임을 다한 a 태그 삭제
}

/**
 * search 버튼 클릭시 동작
 * back-end로 요청 플랫폼/uri를 json화 하여 post
 * 걍 여기서 체크박스에 베스트 퀄리티 클릭하면 바로 
 * 커맨드 꽂아버리고 아니면 메타데이터로 분기하는게 나을듯?
 */
function btnSearch() {
    var obj = new Object();
    obj.platform = document.getElementById("inputGroupSelect").value;
    obj.uri = document.getElementById("inputUrl").value;
    document.getElementById("progressBar").style.display = "";

    isBestQuality = document.getElementById("bestQualityCheckBox").checked
    if (isBestQuality) {
        handleFileDownload('/api/v1/media', obj);
        return;
    }
    targetUri = '/api/v1/meta'

    httpRequest = new XMLHttpRequest();
    httpRequest.onreadystatechange = () => {
        /* readyState가 Done이고 응답 값이 200일 때, 받아온 response로 name과 age를 그려줌 */
        if (httpRequest.readyState === XMLHttpRequest.DONE) {
            var result = httpRequest.response;
            if (httpRequest.status === 200) {
                response = httpRequest.response;
                serveRequest(result)
            } else {
                alert(JSON.stringify(result));
            }
        }
    };
    httpRequest.open('POST', targetUri, true);
    httpRequest.responseType = "json";
    httpRequest.setRequestHeader('Content-Type', 'application/json');
    httpRequest.send(JSON.stringify(obj));
}

function isEmptyObject(param) {
    return Object.keys(param).length === 0 && param.constructor === Object;
}

function serveRequest(response) {
    resultsDiv = document.createElement("div")
    resultsDiv.setAttribute("id", "results")
    resultsDiv.append(createSearchBtn())

    toggleButton = createToggleButton("multiCollapseExample1", "multiCollapseExample2")
    resultsDiv.append(toggleButton)
    resultsDiv.append(createDownloadBtn())

    rowDiv = document.createElement("div")
    rowDiv.setAttribute("class", "row")

    colDiv = createColDiv()
    multiCollapseDiv1 = createMultiCollapseDiv("multiCollapseExample1")

    cardDiv = createCardDiv()

    videoRadioDiv = document.createElement("div")
    videoRadioDiv.setAttribute("class", "my-3")
    for (var key in response.video) {
        radioButton = createRadioButton("video", key)
        label = createLabel("video", key, response)

        cardDiv.append(radioButton, label)
    }
    multiCollapseDiv1.append(cardDiv)

    colDiv.append(multiCollapseDiv1)
    rowDiv.append(colDiv)
    console.log(response);
    if (!isEmptyObject(response.audio)) {
        colDiv = createColDiv()
        multiCollapseDiv2 = createMultiCollapseDiv("multiCollapseExample2")

        cardDiv = createCardDiv()
        audioRadioDiv = document.createElement("div")
        audioRadioDiv.setAttribute("class", "mb-3")
        for (var key in response.audio) {
            radioButton = createRadioButton("audio", key)
            label = createLabel("audio", key, response)

            cardDiv.append(radioButton, label)
        }
        multiCollapseDiv2.append(cardDiv)

        colDiv.append(multiCollapseDiv2)
        rowDiv.append(colDiv)
    }

    resultsDiv.append(rowDiv)

    document.getElementById("results").replaceChildren(resultsDiv)
}

function createSearchBtn() {
    button = document.createElement("button")
    button.setAttribute("type", "button")
    button.setAttribute("onClick", "btnSearch()")
    button.setAttribute("class", "btn btn-primary mx-1")
    button.setAttribute("id", "searchButton")
    button.innerHTML = "Search"
    return button
}
function createDownloadBtn() {
    button = document.createElement("button")
    button.setAttribute("type", "button")
    button.setAttribute("onClick", "btnDownload()")
    button.setAttribute("class", "btn btn-success ")
    button.setAttribute("id", "downloadButton")
    button.innerHTML = "Download"
    return button
}

function createDiv() {
    div = document.createElement("div")
    div.setAttribute("class", "")
    return div
}

function clickFunc() {
    $('.collapse').collapse("toggle")
}

function createColDiv() {
    div = document.createElement("div")
    div.setAttribute("class", "col")
    return div
}

function createMultiCollapseDiv(id) {
    div = document.createElement("div")
    div.setAttribute("class", "collapse multi-collapse")
    div.setAttribute("id", id)
    return div
}

function createToggleButton(videoCard, audioCard) {
    button = document.createElement("button")
    button.setAttribute("class", "btn btn-primary mx-5")
    button.setAttribute("type", "button")
    button.setAttribute("data-toggle", "collapse")
    button.setAttribute("data-target", ".multi-collapse")
    button.setAttribute("onClick", "clickFunc()")
    button.setAttribute("aria-expanded", "false")
    button.setAttribute("id", "toggleButton")
    button.setAttribute("aria-controls", videoCard + " " + audioCard)
    button.innerHTML = "Toggle options"
    return button
}

function createRadioButton(mediaType, key) {
    button = document.createElement("input")
    button.setAttribute("type", "radio")
    button.setAttribute("name", mediaType + "Options")
    button.setAttribute("value", key)
    button.setAttribute("id", mediaType + key)
    button.setAttribute("class", "btn-check")
    return button
}

function createLabel(mediaType, key, response) {
    label = document.createElement("label")
    label.setAttribute("for", mediaType + key)
    if (mediaType == "video") {
        label.setAttribute("class", "btn btn-outline-success mx-2 my-1")
        label.innerHTML = response.video[key].format.split("-")[1] + "(" + response.video[key].filesize + ")"
    } else {
        label.setAttribute("class", "btn btn-outline-primary mx-2 my-1")
        label.innerHTML = response.audio[key].format.split("-")[1] + "(" + response.audio[key].filesize + ")"
    }
    return label
}

function createCardDiv() {
    cardDiv = document.createElement("div")
    cardDiv.setAttribute("class", "card card-body my-2")
    return cardDiv
}