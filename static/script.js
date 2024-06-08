const confessionAmountKey = 'confessionAmount'
const confessionDateKey = 'confessionDate'

const confessional = document.getElementById('confessional')
const confessionalText = document.getElementById('confessional-text')
const confessionsAmount = document.getElementById('confessionsAmount')
const submitConfession = document.getElementById('submit-confession')
const wordCounter = document.getElementById('word-counter')
const playMusic = document.getElementById('play-music')
const recentConfessions = document.getElementById('recent-confessions')
const confessionPublic = document.getElementById('confession-public')

// Ewww dependencies
feather.replace();
const player = new Plyr('#player')
dayjs.extend(window.dayjs_plugin_relativeTime)

// Open websocket connection to websocket server
let wsProtocol = window.location.protocol === 'https:' ? 'wss' : 'ws'
let wsUrl = `${wsProtocol}://${window.location.host}/ws`;
console.log("Connecting to ", wsUrl)
const ws = new WebSocket(wsUrl)

// First message is always an array of confessions
// Confession format {"confession":"","date":""}
let firstMessage = true
ws.onmessage = function (event) {
    const data = JSON.parse(event.data)
    if (firstMessage) {
        if (data !== null) {
            for (let confession of data) {
                recentConfessions.appendChild(ConstructRecentConfession(confession.confession, confession.date))
            }
        }

        firstMessage = false
    } else {
        AddConfession(data)
    }
}

// Adds new confession to end of list and at same time makes sure there is only 5 shown at same time
function AddConfession(confession) {
    recentConfessions.prepend(ConstructRecentConfession(confession.confession, confession.date))
    if (recentConfessions.children.length > 5) {
        recentConfessions.children[recentConfessions.children.length - 1].remove()
    }
}

function ConstructRecentConfession(text, date) {
    let element = document.createElement('confession')

    let content = document.createElement('div')
    content.classList.add('content')

    let dateElement = document.createElement('div')
    dateElement.classList.add('date')
    dateElement.innerText = dayjs(date).format('YYYY-MM-DD HH:mm:ss')
    content.appendChild(dateElement)

    let confession = document.createElement('div')
    confession.classList.add('text')
    confession.innerText = text
    content.appendChild(confession)

    element.appendChild(content)

    return element
}

let confessionAmount = localStorage.getItem(confessionAmountKey) || 0
let confessionDate = localStorage.getItem(confessionDateKey)
if (confessionAmount > 0) {
    updateConfessionsAmount()
}

function updateConfessionsAmount() {
    confessionsAmount.classList.remove('hidden')
    confessionsAmount.innerText = `Confessions: ${confessionAmount}, last confessed: ${dayjs(confessionDate).fromNow()}`
}

confessional.confession.addEventListener('input', (e) => {
    wordCounter.innerText = `Word count: ${e.target.value.trim().split(/\s+/).length-1}`
})

confessional.addEventListener('submit', (e) => {
    e.preventDefault()

    if (confessionalText.value.length < 1) {
        shakeSubmitButton()
        return
    }

    // Submit confession
    fetch('/api/confess', {
        method: 'post',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            'confession': confessionalText.value,
            'public': confessionPublic.checked
        })
    }).then((res) => {
        if (res.ok) {
            // Reset the form
            confessionalText.value = ''
            confessionPublic.checked = false

            // Update the UI
            confessionAmount++
            confessionDate = new Date()
            localStorage.setItem(confessionAmountKey, confessionAmount)
            localStorage.setItem(confessionDateKey, confessionDate)
            updateConfessionsAmount()

            // lol
            confetti()
        } else {
            shakeSubmitButton()
        }
    }).catch((err) => {
        console.log(err)
    })
})

function shakeSubmitButton() {
    submitConfession.classList.remove('shake')
    void submitConfession.offsetWidth; // hacky way to force reflow
    submitConfession.classList.add('shake')
}

playMusic.addEventListener('click', function () {
    player.play()
    playMusic.classList.remove("clickMe")
    playMusic.classList.add('fadeOut')
})