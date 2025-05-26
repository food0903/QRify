"use client";
import Link from "next/link";

const ORANGE = "#FF9900";
const DARK_BG = "#181818";

export default function ExpirationPage() {
    return (
        <div
            style={{
                minHeight: "100vh",
                background: `linear-gradient(135deg, ${ORANGE} 0%, ${DARK_BG} 100%)`,
                color: "#fff",
                fontFamily: "Inter, sans-serif",
                display: "flex",
                flexDirection: "column",
                alignItems: "center",
                justifyContent: "center",
                padding: 0,
            }}
        >
            <header
                style={{
                    padding: "2rem 0 1rem 0",
                    display: "flex",
                    justifyContent: "center",
                    alignItems: "center",
                    width: "100%",
                }}
            >
                <img
                    src="/QRify.png"
                    alt="QRify"
                    style={{
                        width: "120px",
                        height: "120px",
                        maxWidth: "90vw",
                        objectFit: "contain",
                        display: "block",
                    }}
                />
            </header>
            <main
                style={{
                    background: "#222",
                    borderRadius: "18px",
                    boxShadow: "0 4px 32px #0008",
                    width: "min(420px, 95vw)",
                    padding: "2.5rem 2rem",
                    display: "flex",
                    flexDirection: "column",
                    alignItems: "center",
                }}
            >
                <h1 style={{ color: ORANGE, fontWeight: 800, fontSize: "2rem", marginBottom: "1.2rem", textAlign: "center" }}>
                    Your QR code has expired
                </h1>
                <p style={{ color: "#bbb", fontSize: "1.1rem", textAlign: "center", marginBottom: "2rem" }}>
                    This QR code is no longer valid. Please generate a new QR code if you need continued access.
                </p>
                <Link
                    href="/"
                    style={{
                        background: ORANGE,
                        color: DARK_BG,
                        border: "none",
                        borderRadius: "8px",
                        padding: "0.8em 2em",
                        fontWeight: 700,
                        fontSize: "1.1rem",
                        textDecoration: "none",
                        cursor: "pointer",
                        boxShadow: "0 2px 8px #0004",
                        transition: "background 0.2s, color 0.2s",
                    }}
                >
                    Go to Home
                </Link>
            </main>
            <footer style={{ color: "#bbb", fontSize: "0.95rem", marginTop: "auto", padding: "1rem" }}>
                &copy; {new Date().getFullYear()} QRify
            </footer>
        </div>
    );
}
