# BINGO!
> Let's play BINGO!

## 💻 Requirements
- Docker

## 📝 TODO
- [x] Create a new player endpoint
- [x] Raffle a new number

## 🕹️ How to play
- Run:
```
  make server
```

- Open [client.html](cmd/web/client.html) in your browser
  - Press **F12** to see your number

- Raffle a new number calling [host](:8081/bingo/next) until your number appears

- When your number is raffle you was receive a message in browser: `you win`
