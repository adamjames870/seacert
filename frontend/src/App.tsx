import { useEffect, useState } from "react";
import "./App.css";

interface Certificate {
    id: string;
    "cert-type-name": string;
    "cert-number": string;
    "issuer-name": string;
    "issued-date": string;
}

// Format date as dd-MMM-yyyy
function formatDate(dateString: string): string {
    const date = new Date(dateString);
    if (isNaN(date.getTime())) return dateString;
    return date.toLocaleDateString("en-GB", {
        day: "2-digit",
        month: "short",
        year: "numeric",
    }).replace(/ /g, "-");
}

type SortKey = keyof Certificate;

function App() {
    const [certificates, setCertificates] = useState<Certificate[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [sortKey, setSortKey] = useState<SortKey | null>(null);
    const [sortAsc, setSortAsc] = useState(true);

    const fetchCertificates = async () => {
        setLoading(true);
        setError(null);
        try {
            const res = await fetch("http://localhost:8080/api/certificates");
            if (!res.ok) throw new Error(`HTTP error: ${res.status}`);
            const data: Certificate[] = await res.json();
            setCertificates(data);
        } catch (err: unknown) {
            if (err instanceof Error) setError(err.message);
            else setError(String(err));
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchCertificates();
    }, []);

    const handleSort = (key: SortKey) => {
        if (sortKey === key) {
            setSortAsc(!sortAsc); // toggle ascending/descending
        } else {
            setSortKey(key);
            setSortAsc(true);
        }
    };

    const sortedCertificates = [...certificates].sort((a, b) => {
        if (!sortKey) return 0;

        if (sortKey === "issued-date") {
            const aTime = new Date(a[sortKey]).getTime();
            const bTime = new Date(b[sortKey]).getTime();
            return sortAsc ? aTime - bTime : bTime - aTime;
        }

        const aVal = a[sortKey] as string;
        const bVal = b[sortKey] as string;
        return sortAsc ? aVal.localeCompare(bVal) : bVal.localeCompare(aVal);
    });

    if (loading) return <p className="loading">Loading certificates...</p>;
    if (error) return <p className="error">Error: {error}</p>;

    return (
        <div className="app-wrapper">
        <div className="app-container">
            <h1>Certificates</h1>
            <button className="refresh-btn" onClick={fetchCertificates}>
                Refresh
            </button>
            <table className="cert-table">
                <thead>
                <tr>
                    <th onClick={() => handleSort("cert-type-name")}>
                        Cert Type {sortKey === "cert-type-name" ? (sortAsc ? "↑" : "↓") : ""}
                    </th>
                    <th onClick={() => handleSort("cert-number")}>
                        Cert Number {sortKey === "cert-number" ? (sortAsc ? "↑" : "↓") : ""}
                    </th>
                    <th onClick={() => handleSort("issuer-name")}>
                        Issuer {sortKey === "issuer-name" ? (sortAsc ? "↑" : "↓") : ""}
                    </th>
                    <th onClick={() => handleSort("issued-date")}>
                        Issued Date {sortKey === "issued-date" ? (sortAsc ? "↑" : "↓") : ""}
                    </th>
                </tr>
                </thead>
                <tbody>
                {sortedCertificates.map((cert) => (
                    <tr key={cert.id}>
                        <td>{cert["cert-type-name"]}</td>
                        <td>{cert["cert-number"]}</td>
                        <td>{cert["issuer-name"]}</td>
                        <td>{formatDate(cert["issued-date"])}</td>
                    </tr>
                ))}
                </tbody>
            </table>
        </div>
        </div>
    );
}

export default App;
