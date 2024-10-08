const env = process.env.NODE_ENV;

export const BASE_WEBAPP_URL =
  env === "production"
    ? "https://planningpokergame.herokuapp.com"
    : "http://localhost:3000";
export const BASE_URL =
  env === "production"
    ? "https://planningpoker-server.herokuapp.com/api"
    : "http://localhost:8080/api";
export const BASE_WS_URL =
  env === "production"
    ? "wss://planningpoker-server.herokuapp.com/api"
    : "ws://localhost:8080/api";

export const localConstants = {
  ROOM_NAME_KEY: "roomName",
  USERNAME_KEY: "username",
};

export const cookieConstants = {
  USER_KEY: "user",
};

export const ROOM_NAME_INPUT_LIMIT = 30;
export const NAME_INPUT_LIMIT = 15;

export const messages = {
  ROOM_JOINED: "room-joined",
  VOTE_SUBMITTED: "vote-submitted",
  REVEAL_CARDS: "reveal-cards",
  CARDS_REVEALED: "cards-revealed",
  START_NEW_VOTING: "start-new-voting",
  NEW_VOTING_STARTED: "new-voting-started",
  IS_AFK: "is-afk",
  DISCONNECTED: "disconnected",
  KICK_USER: "kick",
  USER_KICKED: "client-kicked",
};

export const gameStates = {
  IN_PROGRESS: "IN_PROGRESS",
  CARDS_REVEALED: "CARDS_REVEALED",
};

export const voteCardValues = {
  CONFUSED: -1, // o card '?'
  NOT_SELECTED: -2,
  PRIVATE: -3, // quando a carta é selecionada mas não revelada, o valor é privado
  EMPTY: -4, // usado somente no lado do cliente, ao reiniciar a sessão de votação
};
