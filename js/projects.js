const container = document.querySelector("[data-projects]");

async function loadProjects() {
    try {
        const res = await fetch("/projects.json");
        const repos = await res.json();

        repos.forEach((repo) => {
            const card = createProjectCard(repo);
            container.appendChild(card);
        });
    } catch (err) {
        console.error("Failed to load projects:", err);
    }
}

function createProjectCard(repo) {
    const card = document.createElement("div");
    card.className = "project-card";
    const title = document.createElement("a");
    title.className = "project-title";
    title.href = repo.html_url;
    title.textContent = repo.name;
    title.target = "_blank";
    const desc = document.createElement("p");
    desc.className = "project-desc";
    desc.textContent = repo.description || "No description";
    const meta = document.createElement("p");
    meta.className = "project-meta";
    meta.textContent = "Updated: " + formatDate(repo.pushed_at);
    const langBar = createLanguageBar(repo.languages);

    card.appendChild(title);
    card.appendChild(meta);
    card.appendChild(desc);
    card.appendChild(langBar);

    return card;
}

function createLanguageBar(languages = {}) {
    const wrapper = document.createElement("div");
    wrapper.className = "lang-wrapper";

    const total = Object.values(languages).reduce((a, b) => a + b, 0);

    const sorted = Object.entries(languages).sort((a, b) => b[1] - a[1]);

    const labels = document.createElement("div");
    labels.className = "lang-labels";

    sorted.forEach(([lang, bytes], index) => {
        const percent = total ? ((bytes / total) * 100).toFixed(1) : 0;

        const span = document.createElement("span");
        span.textContent = `${lang} ${percent}%`;

        if (index < sorted.length - 1) {
            span.textContent += " • ";
        }

        labels.appendChild(span);
    });

    const bar = document.createElement("div");
    bar.className = "lang-bar";

    sorted.forEach(([lang, bytes]) => {
        const segment = document.createElement("div");
        segment.className = "lang-segment";

        segment.style.display = "inline-block";
        segment.style.height = "6px";
        segment.style.width = total ? (bytes / total) * 100 + "%" : "0%";
        segment.style.background = getLangColor(lang);

        bar.appendChild(segment);
    });

    wrapper.appendChild(labels);
    wrapper.appendChild(bar);

    return wrapper;
}

function getLangColor(lang) {
    const colors = {
        Go: "#00ADD8",
        JavaScript: "#f1e05a",
        HTML: "#e34c26",
        CSS: "#563d7c",
        Shell: "#89e051",
        Lua: "#000080",
        C: "#555555",
        Perl: "#0298c3",
        Awk: "#c30e0e",
        Makefile: "#427819",
    };

    return colors[lang] || "#888";
}

function formatDate(dateStr) {
    const date = new Date(dateStr);
    return date.toLocaleDateString();
}

loadProjects();
