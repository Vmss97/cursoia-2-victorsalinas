# Inventory Management System

A full-stack application for managing product inventory, built with Go and React.

## 🚀 Tech Stack

- **Backend**: Go (Golang) 1.22+
  - Concurrent CSV parsing
  - REST API serving JSON
- **Frontend**: React + TypeScript + Vite
  - Tailwind CSS for styling
  - Real-time filtering and status indicators

## 🛠️ Setup & Running

### Backend
The backend serves the data from `inventory.csv` on port `8080`.

1.  Make sure you have Go installed.
2.  Run the service:
    ```bash
    go run main.go
    ```
3.  API: `http://localhost:8080/api/inventory`

### Frontend
The frontend is a modern dashboard running on port `5173`.

1.  Navigate to the frontend folder:
    ```bash
    cd frontend
    ```
2.  Install dependencies:
    ```bash
    npm install
    ```
3.  Start the dev server:
    ```bash
    npm run dev
    ```
4.  Open [http://localhost:5173](http://localhost:5173) in your browser.

## ✨ Features

- **High-Performance Loading**: Uses Go goroutines to parse large CSV files efficiently.
- **Dynamic UI**: Filter products by category instantly.
- **Visual Alerts**: Rows are highlighted in red when stock is 0.
