"use client";
import { useState, useEffect } from "react";
import Link from "next/link";

const ORANGE = "#FF9900";
const DARK_BG = "#181818";

export default function AnalyticsPage() {
    const [qrId, setQrId] = useState("");
    const [scanCount, setScanCount] = useState<number | null>(null);
    const [expiresAt, setExpiresAt] = useState<string | null>(null);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [showModal, setShowModal] = useState(false);
    const [isMobile, setIsMobile] = useState(false);

    useEffect(() => {
        const checkMobile = () => setIsMobile(window.innerWidth < 700);
        checkMobile();
        window.addEventListener("resize", checkMobile);
        return () => window.removeEventListener("resize", checkMobile);
    }, []);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setLoading(true);
        setError(null);
        setScanCount(null);
        setShowModal(false);

        try {
            const res = await fetch(`api/v1/qr/${qrId}/scans`);
            if (!res.ok) {
                let errorMsg = "Failed to fetch scan count.";
                try {
                    const errData = await res.json();
                    if (errData && errData.error) {
                        errorMsg = errData.error;
                    }
                } catch { }
                throw new Error(errorMsg);
            }

            const data = await res.json();
            setScanCount(data.scan_count);
            setExpiresAt(data.expires_at);
            setShowModal(true);
            setError(null);

        } catch (err: any) {
            setError(err.message || "Failed to fetch analytics.");
        } finally {
            setLoading(false);
        }
    };

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
                padding: "0",
            }}
        >
            <Link
                href="/"
                style={{
                    position: "absolute",
                    top: 24,
                    left: 36,
                    color: ORANGE,
                    background: "#181818",
                    fontWeight: 700,
                    textDecoration: "none",
                    fontSize: "1.1rem",
                    display: "flex",
                    alignItems: "center",
                    gap: "0.3em",
                    zIndex: 10,
                    borderRadius: "8px",
                    padding: "0.4em 1em",
                    boxShadow: "0 2px 8px #0004",
                    border: "none",
                    transition: "background 0.2s, color 0.2s",
                }}
                onMouseOver={e => {
                    (e.currentTarget as HTMLElement).style.background = ORANGE;
                    (e.currentTarget as HTMLElement).style.color = "#181818";
                }}
                onMouseOut={e => {
                    (e.currentTarget as HTMLElement).style.background = "#181818";
                    (e.currentTarget as HTMLElement).style.color = ORANGE;
                }}
            >
                ← Back
            </Link>

            <header
                style={{
                    padding: isMobile ? "1.2rem 0 0.5rem 0" : "2rem 0 1rem 0",
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
                        width: isMobile ? "90px" : "160px",
                        height: isMobile ? "90px" : "160px",
                        maxWidth: "90vw",
                        maxHeight: isMobile ? "90px" : "160px",
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
                    width: isMobile ? "98vw" : "min(420px, 95vw)",
                    padding: isMobile ? "1.5rem 1rem" : "2.5rem 2rem",
                    display: "flex",
                    flexDirection: "column",
                    alignItems: "center",
                }}
            >
                <h2 style={{ color: ORANGE, fontWeight: 800, fontSize: "1.3rem", marginBottom: "1.2rem" }}>
                    Check your QR code analytics
                </h2>

                {/* Tutorial Image */}
                <div style={{
                    marginBottom: "1.5rem",
                    width: "100%",
                    display: "flex",
                    justifyContent: "center"
                }}>
                    <img
                        src="/tutorial.png"
                        alt="How to find QR code ID"
                        style={{
                            maxWidth: "100%",
                            height: "auto",
                            borderRadius: "8px",
                            boxShadow: "0 2px 8px #0004"
                        }}
                    />
                </div>

                <div style={{ color: "#bbb", fontSize: "1rem", marginBottom: "1.5rem", textAlign: "center" }}>
                    Enter your QR code ID below.<br />
                    <span style={{ color: ORANGE, fontWeight: 600, opacity: 0.5 }}>
                        (You can find your QR code ID below the generated QR code on the main page or in the downloaded file name)
                    </span>
                </div>
                <form
                    onSubmit={handleSubmit}
                    style={{ display: "flex", gap: "0.5rem", width: "100%", marginBottom: "1.5rem" }}
                >
                    <input
                        type="text"
                        placeholder="Enter your QR code ID"
                        value={qrId}
                        onChange={e => setQrId(e.target.value)}
                        style={{
                            flex: 1,
                            padding: "0.75rem 1rem",
                            borderRadius: "8px",
                            border: "none",
                            fontSize: "1.1rem",
                            background: "#222",
                            color: "#fff",
                            outline: `2px solid ${ORANGE}`,
                        }}
                    />
                    <button
                        type="submit"
                        style={{
                            background: ORANGE,
                            color: "#222",
                            border: "none",
                            borderRadius: "8px",
                            padding: "0.75rem 1.2rem",
                            fontWeight: 700,
                            fontSize: "1.1rem",
                            cursor: "pointer",
                            transition: "background 0.2s",
                            opacity: loading ? 0.7 : 1,
                        }}
                        disabled={loading || !qrId.trim()}
                    >
                        {loading ? "Loading..." : "Check"}
                    </button>
                </form>
                {error && (
                    <div style={{
                        color: "#ff4444",
                        marginBottom: "1rem",
                        textAlign: "center",
                        fontSize: "1rem",
                        padding: "0.5rem 1rem",
                        background: "rgba(255, 68, 68, 0.1)",
                        borderRadius: "8px",
                        width: "100%"
                    }}>
                        {error}
                    </div>
                )}
            </main>

            {showModal && (
                <div
                    style={{
                        position: "fixed",
                        top: 0,
                        left: 0,
                        right: 0,
                        bottom: 0,
                        background: "rgba(0, 0, 0, 0.7)",
                        display: "flex",
                        alignItems: "center",
                        justifyContent: "center",
                        zIndex: 1000,
                    }}
                    onClick={() => setShowModal(false)}
                >
                    <div
                        style={{
                            background: "#222",
                            borderRadius: "18px",
                            boxShadow: "0 4px 32px #0008",
                            width: "min(420px, 95vw)",
                            padding: "2rem",
                            display: "flex",
                            flexDirection: "column",
                            alignItems: "center",
                            gap: "1rem",
                            position: "relative",
                        }}
                        onClick={e => e.stopPropagation()}
                    >
                        <button
                            onClick={() => setShowModal(false)}
                            style={{
                                position: "absolute",
                                top: "1rem",
                                right: "1rem",
                                background: "none",
                                border: "none",
                                color: "#666",
                                fontSize: "1.5rem",
                                cursor: "pointer",
                                padding: "0.5rem",
                                display: "flex",
                                alignItems: "center",
                                justifyContent: "center",
                                transition: "color 0.2s",
                            }}
                            onMouseOver={e => {
                                (e.currentTarget as HTMLElement).style.color = ORANGE;
                            }}
                            onMouseOut={e => {
                                (e.currentTarget as HTMLElement).style.color = "#666";
                            }}
                        >
                            ×
                        </button>

                        {error ? (
                            <div style={{
                                color: "#ff4444",
                                textAlign: "center",
                                fontSize: "1.1rem",
                                padding: "0.5rem 1rem",
                                background: "rgba(255, 68, 68, 0.1)",
                                borderRadius: "8px",
                                width: "100%"
                            }}>
                                {error}
                            </div>
                        ) : (
                            <>
                                <div style={{
                                    fontSize: "1.3rem",
                                    fontWeight: 700,
                                    color: ORANGE,
                                    textAlign: "center"
                                }}>
                                    Scan Statistics
                                </div>
                                <div style={{
                                    fontSize: "2.5rem",
                                    fontWeight: 800,
                                    color: "#fff",
                                    textAlign: "center"
                                }}>
                                    {scanCount}
                                </div>
                                <div style={{
                                    fontSize: "1.1rem",
                                    color: "#bbb",
                                    textAlign: "center"
                                }}>
                                    {scanCount === 1 ? "scan" : "scans"} recorded
                                </div>
                                {expiresAt && (
                                    <div style={{ fontSize: "1.1rem", color: "red", textAlign: "center" }}>
                                        {(() => {
                                            const exp = new Date(expiresAt);
                                            if (
                                                exp.getUTCFullYear() === 1 &&
                                                exp.getUTCMonth() === 0 &&
                                                exp.getUTCDate() === 1 &&
                                                exp.getUTCHours() === 0 &&
                                                exp.getUTCMinutes() === 0 &&
                                                exp.getUTCSeconds() === 0
                                            ) {
                                                return "This QR code never expires";
                                            }
                                            const now = new Date();
                                            if (exp < now) {
                                                return "QR code has expired";
                                            }
                                            const diff = exp.getTime() - now.getTime();
                                            const days = Math.floor(diff / (1000 * 60 * 60 * 24));
                                            const hours = Math.floor((diff / (1000 * 60 * 60)) % 24);
                                            const minutes = Math.floor((diff / (1000 * 60)) % 60);
                                            return `QR code is valid for: ${days}d ${hours}h ${minutes}m`;
                                        })()}
                                    </div>
                                )}
                            </>
                        )}
                    </div>
                </div>
            )}

            <footer style={{ color: "#bbb", fontSize: "0.95rem", marginTop: "auto", padding: "1rem" }}>
                &copy; {new Date().getFullYear()} QRify
            </footer>
        </div>
    );
}
