"use client";
import { useState, useEffect } from "react";
import Link from "next/link";

const ORANGE = "#FF9900";
const DARK_BG = "#181818";
const apiUrl = process.env.NEXT_PUBLIC_API_URL;

const EXPIRATION_OPTIONS = [
  { label: "12 hours", value: 12 * 60 * 60 },
  { label: "1 day", value: 24 * 60 * 60 },
  { label: "30 days", value: 30 * 24 * 60 * 60 },
  { label: "90 days", value: 90 * 24 * 60 * 60 },
  { label: "180 days", value: 180 * 24 * 60 * 60 },
  { label: "1 year", value: 365 * 24 * 60 * 60 },
  { label: "5 years", value: 5 * 365 * 24 * 60 * 60 },
  { label: "Never", value: 0 }
];

export default function Home() {
  const [url, setUrl] = useState("");
  const [qrImage, setQrImage] = useState<string | null>(null);
  const [qrId, setQrId] = useState<string | null>(null);
  const [copied, setCopied] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const [isMobile, setIsMobile] = useState(false);
  const [expirationTime, setExpirationTime] = useState<string>("");

  useEffect(() => {
    const checkMobile = () => setIsMobile(window.innerWidth < 700);
    checkMobile();
    window.addEventListener("resize", checkMobile);
    return () => window.removeEventListener("resize", checkMobile);
  }, []);

  const handleInput = (e: React.ChangeEvent<HTMLInputElement>) => {
    setUrl(e.target.value);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    try {
      const normalizedUrl = url.trim()
        .replace(/^https?:\/\/(www\.)?/, 'https://')
        .replace(/\/+$/, '');

      const searchUrl = `${apiUrl}/v1/qr?url=${encodeURIComponent(normalizedUrl)}`;
      const searchRes = await fetch(searchUrl);

      let qrObj;

      if (searchRes.ok) {
        const data = await searchRes.json();
        const isNeverExpires = (() => {
          const exp = new Date(data.expires_at);
          return (
            exp.getUTCFullYear() === 1 &&
            exp.getUTCMonth() === 0 &&
            exp.getUTCDate() === 1 &&
            exp.getUTCHours() === 0 &&
            exp.getUTCMinutes() === 0 &&
            exp.getUTCSeconds() === 0
          );
        })();

        if (
          data.expires_at &&
          !isNeverExpires &&
          new Date(data.expires_at).getTime() < Date.now()
        ) {
          const createRes = await fetch(`${apiUrl}/v1/qr`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ url: normalizedUrl, expires_in_sec: Number(expirationTime) }),
          });
          if (!createRes.ok) {
            let errorMsg = "Failed to create QR code.";
            try {
              const errData = await createRes.json();
              if (errData && errData.error) {
                errorMsg = errData.error;
              }
            } catch { }
            throw new Error(errorMsg);
          }
          qrObj = await createRes.json();
        } else {
          qrObj = data;
        }
      } else if (searchRes.status === 404) {
        const createRes = await fetch(`${apiUrl}/v1/qr`, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ url: normalizedUrl, expires_in_sec: Number(expirationTime) }),
        });
        if (!createRes.ok) {
          let errorMsg = "Failed to create QR code.";
          try {
            const errData = await createRes.json();
            if (errData && errData.error) {
              errorMsg = errData.error;
            }
          } catch { }
          throw new Error(errorMsg);
        }
        qrObj = await createRes.json();
      } else {
        let errorMsg = "Failed to search QR code.";
        try {
          const errData = await searchRes.json();
          if (errData && errData.error) {
            errorMsg = errData.error;
          }
        } catch { }
        throw new Error(errorMsg);
      }

      setQrId(qrObj.id);
      setQrImage(`data:image/png;base64,${qrObj.image_base64}`);
    } catch (err: any) {
      setError(err.message || "Something went wrong.");
    } finally {
      setLoading(false);
    }
  };

  const isGenerateDisabled = loading || !url.trim() || expirationTime === "";

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
        position: "relative",
      }}
    >

      <Link
        href="/analytics"
        style={{
          position: "absolute",
          top: 24,
          right: 36,
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
        Analytics <span style={{ fontSize: "1.3em" }}>→</span>
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
          display: "flex",
          flexDirection: isMobile ? "column" : "row",
          width: isMobile ? "98vw" : "min(900px, 95vw)",
          minHeight: isMobile ? "auto" : "480px",
          margin: isMobile ? "1rem 0" : "2rem 0",
          overflow: "hidden",
        }}
      >
        {/* Left: Input */}
        <section
          style={{
            flex: 1,
            padding: isMobile ? "1.5rem 1rem" : "2.5rem 2rem",
            display: "flex",
            flexDirection: "column",
            justifyContent: "center",
            borderRight: isMobile ? "none" : "1px solid #333",
            borderBottom: isMobile ? "1px solid #333" : "none",
            background: "#181818",
          }}
        >
          <h1 style={{
            fontSize: isMobile ? "1.4rem" : "2.2rem",
            fontWeight: 800,
            marginBottom: isMobile ? "1rem" : "1.5rem",
            color: ORANGE
          }}>
            Convert your Link<br />to QR code
          </h1>
          {/* Form with input and submit button */}
          <form
            style={{ display: "flex", gap: "0.5rem", marginBottom: "1.2rem" }}
            onSubmit={handleSubmit}
          >
            <input
              type="text"
              placeholder="Enter or paste url"
              value={url}
              onChange={handleInput}
              style={{
                flex: 1,
                padding: "0.75rem 1rem",
                borderRadius: "8px",
                border: "none",
                fontSize: "1.1rem",
                background: "#222",
                color: "#fff",
                outline: "2px solid #ff9900",
              }}
            />
            <select
              value={expirationTime}
              onChange={e => setExpirationTime(e.target.value)}
              style={{
                padding: "0.75rem 1rem",
                borderRadius: "8px",
                border: "none",
                fontSize: "1.1rem",
                background: "#222",
                color: "#fff",
                outline: "2px solid #FF9900",
                cursor: "pointer",
                minWidth: "120px",
                marginLeft: "0.5rem"
              }}
            >
              <option value="" disabled>
                Expiration
              </option>
              {EXPIRATION_OPTIONS.map(option => (
                <option key={option.value} value={option.value}>
                  {option.label}
                </option>
              ))}
            </select>
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
                cursor: isGenerateDisabled ? "not-allowed" : "pointer",
                opacity: isGenerateDisabled ? 0.5 : 1,
                transition: "background 0.2s, opacity 0.2s",
              }}
              disabled={isGenerateDisabled}
            >
              {loading ? "Loading..." : "Generate"}
            </button>
          </form>
          <div style={{ fontSize: "0.95rem", color: "#bbb" }}>
            Your QR code will be generated automatically. <br />
            The generated QR code will open this URL.
          </div>
          {error && (
            <div
              style={{
                color: "#ff4444",
                marginTop: "1rem",
                textAlign: "center",
                fontSize: "1rem",
                padding: "0.5rem 1rem",
                background: "rgba(255, 68, 68, 0.1)",
                borderRadius: "8px",
                width: "100%"
              }}
            >
              {error}
            </div>
          )}
        </section>

        {/* Right: QR code and ID */}
        <section
          style={{
            flex: 1,
            padding: isMobile ? "1.5rem 1rem" : "2.5rem 2rem",
            display: "flex",
            flexDirection: "column",
            alignItems: "center",
            justifyContent: "center",
            background: "#232323",
          }}
        >
          <div
            style={{
              width: isMobile ? 160 : 220,
              height: isMobile ? 160 : 220,
              background: "#333",
              borderRadius: "12px",
              display: "flex",
              alignItems: "center",
              justifyContent: "center",
              marginBottom: isMobile ? "1rem" : "1.5rem",
            }}
          >
            {qrImage ? (
              <img src={qrImage} alt="QR code" style={{ width: isMobile ? 140 : 200, height: isMobile ? 140 : 200 }} />
            ) : (
              <div style={{ color: "#555", fontSize: isMobile ? "1rem" : "1.2rem" }}>QR code preview</div>
            )}
          </div>
          <div
            style={{
              marginBottom: "0.5rem",
              display: "flex",
              flexDirection: "column",
              alignItems: "center"
            }}
          >
            <button
              type="button"
              disabled={!qrId}
              onClick={() => {
                if (qrId) {
                  navigator.clipboard.writeText(qrId);
                  setCopied(true);
                  setTimeout(() => setCopied(false), 1200);
                }
              }}
              style={{
                background: "none",
                border: "none",
                color: qrId ? ORANGE : "#bbb",
                fontWeight: 700,
                fontSize: "1.2rem",
                letterSpacing: "0.04em",
                opacity: qrId ? 1 : 0.5,
                cursor: qrId ? "pointer" : "not-allowed",
                outline: "none",
                padding: 0,
                transition: "opacity 0.2s",
                display: "flex",
                alignItems: "center",
                gap: "0.5em",
              }}
              aria-label={qrId ? "Copy QR code ID" : "QR code ID will appear here"}
            >
              {qrId ? (
                <>
                  {qrId}
                  <span style={{ fontSize: "1.2em" }}>📋</span>
                </>
              ) : (
                <span>your qr code id will appear here</span>
              )}
            </button>
            {copied && (
              <span
                style={{
                  marginTop: 6,
                  fontSize: "1em",
                  color: "white",

                  padding: "0.1em 0.7em",
                  fontWeight: 600,
                  letterSpacing: "0.03em",
                  boxShadow: "0 2px 8px #0002"
                }}
              >
                Copied!
              </span>
            )}
          </div>
          <button
            type="button"
            disabled={!qrImage}
            onClick={() => {
              if (qrImage) {
                const link = document.createElement("a");
                link.href = qrImage;
                link.download = `${qrId}.png`;
                link.click();
              }
            }}
            style={{
              marginTop: "1.2rem",
              background: ORANGE,
              color: "#181818",
              border: "none",
              borderRadius: "8px",
              padding: "0.7em 1.5em",
              fontWeight: 700,
              fontSize: "1.1rem",
              cursor: qrImage ? "pointer" : "not-allowed",
              opacity: qrImage ? 1 : 0.5,
              boxShadow: "0 2px 8px #0004",
              transition: "opacity 0.2s",
            }}
          >
            Download
          </button>
        </section>
      </main>
      <footer style={{ color: "#bbb", fontSize: "0.95rem", marginTop: "auto", padding: "1rem" }}>
        &copy; {new Date().getFullYear()} QRify
      </footer>
    </div>
  );
}