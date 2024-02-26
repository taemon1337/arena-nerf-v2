import { writable } from 'svelte/store'

export const API = import.meta.env.VITE_API || "/api/v1"

export const api = (path) => {
  return fetch(API + path, { credentials: 'include' })
}

export const sendaction = (action, payload) => {
  return fetch(API + "/do/"+encodeURIComponent(action), {
    method: "POST",
    credentials: 'include',
    body: JSON.stringify(payload)
  })
}

export const uuid = writable("current")
export const scoreboard = writable({})
export const nodeboard = writable({})
export const nodes = writable([])
export const teams = writable([])
export const gameevents = writable([])
export const gamelist = writable([])

export const currentGame = writable({
  start_at: "",
  end_at: "",
  length: "",
  completed: false,
  status: "no active game",
  winner: "",
  highscore: 0,
})

export async function fetchGame(id) {
  const res = await api("/games/" + id)
  let data = await res.json()
  if (data.stats) {
    currentGame.update(() => Object.assign({start_at: data.stats.start_at, end_at: data.stats.end_at, length: data.stats.length, completed: data.stats.completed, status: data.stats.status, winner: data.stats.winner, highscore: data.stats.highscore }))
    scoreboard.update(() => data.stats.scoreboard)
    nodeboard.update(() => data.stats.nodeboard)
    nodes.update(() => data.stats.nodes)
    teams.update(() => data.stats.teams)
    gameevents.update(() => data.stats.events)
  } else if (data.games) {
    gamelist.update(() => data.games)
  } else {
    console.error("Unexpected API response, expected a 'stats' key.", data)
  }
}

export async function fetchGames() {
  return fetchGame("all")
}

let poller;
let pollgames;

export const pollGame = function(id) {
  if (poller) {
    clearInterval(poller)
  }
  poller = setInterval(function() {
    fetchGame(id)
  }, 5000)
}

export const pollGames = function() {
  if (pollgames) {
    clearInterval(pollgames)
  }
  pollgames = setInterval(function() {
    fetchGame("all")
  }, 5000)
}
