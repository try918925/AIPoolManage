package main

const webUIHTML = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>AI Gateway · 统一模型管理平台</title>
<link href="https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700;800&display=swap" rel="stylesheet">
<style>
*,*::before,*::after{margin:0;padding:0;box-sizing:border-box}

/* ═══ Theme: Midnight (default) ═══ */
:root,[data-theme="midnight"]{
  --bg-primary:#050509;--bg-secondary:#0c0c14;--bg-card:#111119;--bg-card-hover:#18182a;--bg-input:#0e0e1a;
  --bg-surface:rgba(255,255,255,0.03);--bg-sidebar:rgba(12,12,20,0.85);--bg-table:rgba(12,12,20,0.4);
  --bg-toast:rgba(17,17,25,0.9);--bg-modal:#111119;--bg-login:rgba(17,17,25,0.8);
  --border:rgba(255,255,255,0.06);--border-light:rgba(255,255,255,0.1);--border-accent:rgba(99,102,241,0.3);
  --text-primary:#f0f0f8;--text-secondary:#8b8ba8;--text-muted:#55556e;
  --accent:#6366f1;--accent-light:#818cf8;--accent-glow:rgba(99,102,241,0.25);--accent-subtle:rgba(99,102,241,0.08);
  --success:#10b981;--success-glow:rgba(16,185,129,0.15);--warning:#f59e0b;--danger:#ef4444;--danger-glow:rgba(239,68,68,0.15);--info:#3b82f6;
  --gradient-1:linear-gradient(135deg,#6366f1 0%,#8b5cf6 100%);
  --gradient-2:linear-gradient(135deg,#10b981 0%,#06b6d4 100%);
  --gradient-3:linear-gradient(135deg,#f43f5e 0%,#ec4899 100%);
  --gradient-4:linear-gradient(135deg,#f59e0b 0%,#ef4444 100%);
  --gradient-bg:linear-gradient(160deg,rgba(99,102,241,0.04) 0%,transparent 40%,transparent 60%,rgba(16,185,129,0.03) 100%);
  --shadow-sm:0 1px 3px rgba(0,0,0,0.4);--shadow-md:0 4px 16px rgba(0,0,0,0.5);--shadow-lg:0 12px 48px rgba(0,0,0,0.6);
  --shadow-glow:0 0 40px rgba(99,102,241,0.15);--shadow-card:0 1px 3px rgba(0,0,0,0.3),0 0 0 1px rgba(255,255,255,0.03);
  --scrollbar-thumb:rgba(255,255,255,0.08);--scrollbar-hover:rgba(255,255,255,0.14);
  --tr-hover:rgba(99,102,241,0.03);--logo-glow:rgba(99,102,241,0.3);--login-glow:rgba(99,102,241,0.08);
  --select-arrow:%238b8ba8;--tag-openai-bg:rgba(16,185,129,0.08);--tag-openai-c:#34d399;--tag-openai-b:rgba(16,185,129,0.12);
  --tag-claude-bg:rgba(217,119,50,0.08);--tag-claude-c:#f0a060;--tag-claude-b:rgba(217,119,50,0.12);
  --tag-qwen-bg:rgba(59,130,246,0.08);--tag-qwen-c:#60a5fa;--tag-qwen-b:rgba(59,130,246,0.12);
  --is-dark:1;
}

/* ═══ Theme: Cyberpunk ═══ */
[data-theme="cyberpunk"]{
  --bg-primary:#09090d;--bg-secondary:#0d0d15;--bg-card:#12111c;--bg-card-hover:#1a1928;--bg-input:#0f0e18;
  --bg-surface:rgba(0,255,200,0.02);--bg-sidebar:rgba(13,13,21,0.88);--bg-table:rgba(13,13,21,0.5);
  --bg-toast:rgba(18,17,28,0.92);--bg-modal:#12111c;--bg-login:rgba(18,17,28,0.85);
  --border:rgba(0,255,200,0.07);--border-light:rgba(0,255,200,0.12);--border-accent:rgba(0,255,200,0.3);
  --text-primary:#e0ffe8;--text-secondary:#78b8a0;--text-muted:#3d6858;
  --accent:#00ffc8;--accent-light:#5cffda;--accent-glow:rgba(0,255,200,0.2);--accent-subtle:rgba(0,255,200,0.06);
  --success:#00ffc8;--success-glow:rgba(0,255,200,0.15);--warning:#ffd000;--danger:#ff3860;--danger-glow:rgba(255,56,96,0.15);--info:#00d4ff;
  --gradient-1:linear-gradient(135deg,#00ffc8 0%,#00d4ff 100%);
  --gradient-2:linear-gradient(135deg,#00ffc8 0%,#7cff6b 100%);
  --gradient-3:linear-gradient(135deg,#ff3860 0%,#ff6090 100%);
  --gradient-4:linear-gradient(135deg,#ffd000 0%,#ff8c00 100%);
  --gradient-bg:linear-gradient(160deg,rgba(0,255,200,0.03) 0%,transparent 40%,transparent 60%,rgba(0,212,255,0.02) 100%);
  --shadow-sm:0 1px 3px rgba(0,0,0,0.5);--shadow-md:0 4px 16px rgba(0,0,0,0.6);--shadow-lg:0 12px 48px rgba(0,0,0,0.7);
  --shadow-glow:0 0 40px rgba(0,255,200,0.1);--shadow-card:0 1px 3px rgba(0,0,0,0.4),0 0 0 1px rgba(0,255,200,0.04);
  --scrollbar-thumb:rgba(0,255,200,0.1);--scrollbar-hover:rgba(0,255,200,0.18);
  --tr-hover:rgba(0,255,200,0.03);--logo-glow:rgba(0,255,200,0.3);--login-glow:rgba(0,255,200,0.06);
  --select-arrow:%2378b8a0;--tag-openai-bg:rgba(0,255,200,0.06);--tag-openai-c:#5cffda;--tag-openai-b:rgba(0,255,200,0.1);
  --tag-claude-bg:rgba(255,150,50,0.06);--tag-claude-c:#ffa060;--tag-claude-b:rgba(255,150,50,0.1);
  --tag-qwen-bg:rgba(0,212,255,0.06);--tag-qwen-c:#60d4ff;--tag-qwen-b:rgba(0,212,255,0.1);
  --is-dark:1;
}
[data-theme="cyberpunk"] .btn-primary{color:#09090d;font-weight:700}
[data-theme="cyberpunk"] .sidebar-logo h1{background:var(--gradient-1);-webkit-background-clip:text;-webkit-text-fill-color:transparent}
[data-theme="cyberpunk"] .login-card h2{background:var(--gradient-1);-webkit-background-clip:text;-webkit-text-fill-color:transparent}

/* ═══ Theme: Aurora ═══ */
[data-theme="aurora"]{
  --bg-primary:#070710;--bg-secondary:#0b0b18;--bg-card:#10101e;--bg-card-hover:#171730;--bg-input:#0d0d1c;
  --bg-surface:rgba(168,85,247,0.03);--bg-sidebar:rgba(11,11,24,0.88);--bg-table:rgba(11,11,24,0.5);
  --bg-toast:rgba(16,16,30,0.92);--bg-modal:#10101e;--bg-login:rgba(16,16,30,0.85);
  --border:rgba(168,85,247,0.08);--border-light:rgba(168,85,247,0.14);--border-accent:rgba(168,85,247,0.3);
  --text-primary:#f0e8ff;--text-secondary:#9b8cc0;--text-muted:#5c4d78;
  --accent:#a855f7;--accent-light:#c084fc;--accent-glow:rgba(168,85,247,0.22);--accent-subtle:rgba(168,85,247,0.07);
  --success:#22d3ee;--success-glow:rgba(34,211,238,0.15);--warning:#fbbf24;--danger:#f43f5e;--danger-glow:rgba(244,63,94,0.15);--info:#60a5fa;
  --gradient-1:linear-gradient(135deg,#a855f7 0%,#ec4899 100%);
  --gradient-2:linear-gradient(135deg,#22d3ee 0%,#a855f7 100%);
  --gradient-3:linear-gradient(135deg,#f43f5e 0%,#fb923c 100%);
  --gradient-4:linear-gradient(135deg,#fbbf24 0%,#f43f5e 100%);
  --gradient-bg:linear-gradient(160deg,rgba(168,85,247,0.04) 0%,transparent 35%,transparent 65%,rgba(236,72,153,0.03) 100%);
  --shadow-sm:0 1px 3px rgba(0,0,0,0.4);--shadow-md:0 4px 16px rgba(0,0,0,0.5);--shadow-lg:0 12px 48px rgba(0,0,0,0.6);
  --shadow-glow:0 0 40px rgba(168,85,247,0.12);--shadow-card:0 1px 3px rgba(0,0,0,0.3),0 0 0 1px rgba(168,85,247,0.04);
  --scrollbar-thumb:rgba(168,85,247,0.1);--scrollbar-hover:rgba(168,85,247,0.18);
  --tr-hover:rgba(168,85,247,0.03);--logo-glow:rgba(168,85,247,0.3);--login-glow:rgba(168,85,247,0.06);
  --select-arrow:%239b8cc0;--tag-openai-bg:rgba(34,211,238,0.06);--tag-openai-c:#67e8f9;--tag-openai-b:rgba(34,211,238,0.1);
  --tag-claude-bg:rgba(251,146,60,0.06);--tag-claude-c:#fdba74;--tag-claude-b:rgba(251,146,60,0.1);
  --tag-qwen-bg:rgba(96,165,250,0.06);--tag-qwen-c:#93bbfd;--tag-qwen-b:rgba(96,165,250,0.1);
  --is-dark:1;
}

/* ═══ Theme: Daylight ═══ */
[data-theme="daylight"]{
  --bg-primary:#f8f9fc;--bg-secondary:#ffffff;--bg-card:#ffffff;--bg-card-hover:#f1f3f9;--bg-input:#f3f4f8;
  --bg-surface:rgba(0,0,0,0.02);--bg-sidebar:rgba(255,255,255,0.92);--bg-table:rgba(255,255,255,0.8);
  --bg-toast:rgba(255,255,255,0.95);--bg-modal:#ffffff;--bg-login:rgba(255,255,255,0.9);
  --border:rgba(0,0,0,0.08);--border-light:rgba(0,0,0,0.12);--border-accent:rgba(99,102,241,0.35);
  --text-primary:#1a1a2e;--text-secondary:#64648a;--text-muted:#9e9eb8;
  --accent:#6366f1;--accent-light:#4f46e5;--accent-glow:rgba(99,102,241,0.18);--accent-subtle:rgba(99,102,241,0.06);
  --success:#059669;--success-glow:rgba(5,150,105,0.12);--warning:#d97706;--danger:#dc2626;--danger-glow:rgba(220,38,38,0.12);--info:#2563eb;
  --gradient-1:linear-gradient(135deg,#6366f1 0%,#8b5cf6 100%);
  --gradient-2:linear-gradient(135deg,#059669 0%,#0891b2 100%);
  --gradient-3:linear-gradient(135deg,#dc2626 0%,#db2777 100%);
  --gradient-4:linear-gradient(135deg,#d97706 0%,#dc2626 100%);
  --gradient-bg:linear-gradient(160deg,rgba(99,102,241,0.03) 0%,transparent 40%,transparent 60%,rgba(5,150,105,0.02) 100%);
  --shadow-sm:0 1px 3px rgba(0,0,0,0.06);--shadow-md:0 4px 16px rgba(0,0,0,0.08);--shadow-lg:0 12px 48px rgba(0,0,0,0.12);
  --shadow-glow:0 0 40px rgba(99,102,241,0.08);--shadow-card:0 1px 4px rgba(0,0,0,0.05),0 0 0 1px rgba(0,0,0,0.04);
  --scrollbar-thumb:rgba(0,0,0,0.12);--scrollbar-hover:rgba(0,0,0,0.2);
  --tr-hover:rgba(99,102,241,0.04);--logo-glow:rgba(99,102,241,0.2);--login-glow:rgba(99,102,241,0.04);
  --select-arrow:%2364648a;--tag-openai-bg:rgba(5,150,105,0.08);--tag-openai-c:#059669;--tag-openai-b:rgba(5,150,105,0.15);
  --tag-claude-bg:rgba(217,119,6,0.08);--tag-claude-c:#b45309;--tag-claude-b:rgba(217,119,6,0.15);
  --tag-qwen-bg:rgba(37,99,235,0.08);--tag-qwen-c:#2563eb;--tag-qwen-b:rgba(37,99,235,0.15);
  --is-dark:0;
}
[data-theme="daylight"] body::before,[data-theme="daylight"] body::after{opacity:0.5}
[data-theme="daylight"] .sidebar-logo h1{background:var(--gradient-1);-webkit-background-clip:text;-webkit-text-fill-color:transparent}
[data-theme="daylight"] .login-card h2{background:var(--gradient-1);-webkit-background-clip:text;-webkit-text-fill-color:transparent}
[data-theme="daylight"] .login-card::before{background:linear-gradient(90deg,transparent,rgba(99,102,241,0.2),transparent)}
[data-theme="daylight"] .badge-success{background:rgba(5,150,105,0.08);border-color:rgba(5,150,105,0.15)}
[data-theme="daylight"] .badge-danger{background:rgba(220,38,38,0.08);border-color:rgba(220,38,38,0.15)}
[data-theme="daylight"] .badge-warning{background:rgba(217,119,6,0.08);border-color:rgba(217,119,6,0.15)}
[data-theme="daylight"] .badge-info{background:rgba(37,99,235,0.08);border-color:rgba(37,99,235,0.15)}
[data-theme="daylight"] .badge-purple{background:rgba(99,102,241,0.08);border-color:rgba(99,102,241,0.15)}

/* ═══ Theme: Sunset ═══ */
[data-theme="sunset"]{
  --bg-primary:#0f0a08;--bg-secondary:#161010;--bg-card:#1c1412;--bg-card-hover:#261c18;--bg-input:#140f0d;
  --bg-surface:rgba(251,146,60,0.03);--bg-sidebar:rgba(22,16,16,0.88);--bg-table:rgba(22,16,16,0.5);
  --bg-toast:rgba(28,20,18,0.92);--bg-modal:#1c1412;--bg-login:rgba(28,20,18,0.85);
  --border:rgba(251,146,60,0.08);--border-light:rgba(251,146,60,0.14);--border-accent:rgba(251,146,60,0.3);
  --text-primary:#fde8d8;--text-secondary:#b89080;--text-muted:#785848;
  --accent:#f97316;--accent-light:#fb923c;--accent-glow:rgba(249,115,22,0.22);--accent-subtle:rgba(249,115,22,0.07);
  --success:#34d399;--success-glow:rgba(52,211,153,0.15);--warning:#fbbf24;--danger:#ef4444;--danger-glow:rgba(239,68,68,0.15);--info:#60a5fa;
  --gradient-1:linear-gradient(135deg,#f97316 0%,#ef4444 100%);
  --gradient-2:linear-gradient(135deg,#34d399 0%,#fbbf24 100%);
  --gradient-3:linear-gradient(135deg,#ef4444 0%,#ec4899 100%);
  --gradient-4:linear-gradient(135deg,#fbbf24 0%,#f97316 100%);
  --gradient-bg:linear-gradient(160deg,rgba(249,115,22,0.04) 0%,transparent 35%,transparent 65%,rgba(239,68,68,0.03) 100%);
  --shadow-sm:0 1px 3px rgba(0,0,0,0.5);--shadow-md:0 4px 16px rgba(0,0,0,0.6);--shadow-lg:0 12px 48px rgba(0,0,0,0.7);
  --shadow-glow:0 0 40px rgba(249,115,22,0.1);--shadow-card:0 1px 3px rgba(0,0,0,0.4),0 0 0 1px rgba(251,146,60,0.04);
  --scrollbar-thumb:rgba(251,146,60,0.1);--scrollbar-hover:rgba(251,146,60,0.18);
  --tr-hover:rgba(249,115,22,0.03);--logo-glow:rgba(249,115,22,0.3);--login-glow:rgba(249,115,22,0.06);
  --select-arrow:%23b89080;--tag-openai-bg:rgba(52,211,153,0.06);--tag-openai-c:#6ee7b7;--tag-openai-b:rgba(52,211,153,0.1);
  --tag-claude-bg:rgba(251,146,60,0.06);--tag-claude-c:#fdba74;--tag-claude-b:rgba(251,146,60,0.1);
  --tag-qwen-bg:rgba(96,165,250,0.06);--tag-qwen-c:#93bbfd;--tag-qwen-b:rgba(96,165,250,0.1);
  --is-dark:1;
}

/* ═══ Theme Switcher Styles ═══ */
.theme-fab{position:fixed;bottom:28px;right:28px;z-index:150;width:46px;height:46px;border-radius:50%;background:var(--gradient-1);border:none;cursor:pointer;display:flex;align-items:center;justify-content:center;font-size:20px;box-shadow:0 4px 20px var(--accent-glow),0 0 0 1px rgba(255,255,255,0.06);transition:all 0.3s cubic-bezier(0.34,1.56,0.64,1);color:#fff}
.theme-fab:hover{transform:scale(1.1);box-shadow:0 6px 28px var(--accent-glow)}
.theme-fab.open{transform:rotate(45deg) scale(1.05)}
.theme-panel{position:fixed;bottom:86px;right:28px;z-index:150;background:var(--bg-card);border:1px solid var(--border-light);border-radius:var(--radius-lg);padding:8px;width:230px;box-shadow:var(--shadow-lg);opacity:0;transform:translateY(12px) scale(0.95);pointer-events:none;transition:all 0.25s cubic-bezier(0.34,1.56,0.64,1)}
.theme-panel.open{opacity:1;transform:translateY(0) scale(1);pointer-events:auto}
.theme-panel-title{font-size:10px;font-weight:700;text-transform:uppercase;letter-spacing:1.8px;color:var(--text-muted);padding:10px 12px 6px}
.theme-option{display:flex;align-items:center;gap:12px;padding:10px 12px;border-radius:var(--radius);cursor:pointer;transition:all 0.15s ease;border:1px solid transparent;margin-bottom:2px}
.theme-option:last-child{margin-bottom:0}
.theme-option:hover{background:var(--accent-subtle);border-color:var(--border)}
.theme-option.active{background:var(--accent-subtle);border-color:var(--border-accent)}
.theme-option-colors{display:flex;gap:3px;flex-shrink:0}
.theme-option-dot{width:14px;height:14px;border-radius:50%;box-shadow:inset 0 1px 2px rgba(0,0,0,0.2)}
.theme-option-info{flex:1;min-width:0}
.theme-option-name{font-size:13px;font-weight:600;color:var(--text-primary);line-height:1.2}
.theme-option-desc{font-size:10px;color:var(--text-muted);margin-top:2px}

/* ═══ Shared Styles ═══ */
:root{
  --radius:10px;
  --radius-lg:14px;
  --radius-xl:20px;
  --transition:all 0.2s cubic-bezier(0.4,0,0.2,1);
  --transition-slow:all 0.4s cubic-bezier(0.4,0,0.2,1);
}
html{font-size:14px;scroll-behavior:smooth}
body{font-family:'Inter',system-ui,-apple-system,sans-serif;background:var(--bg-primary);color:var(--text-primary);min-height:100vh;overflow-x:hidden}
body::before{content:'';position:fixed;top:0;left:0;right:0;bottom:0;background:var(--gradient-bg);pointer-events:none;z-index:0}
body::after{content:'';position:fixed;top:-50%;left:-50%;width:200%;height:200%;background:radial-gradient(circle at 30% 20%,rgba(99,102,241,0.03) 0%,transparent 40%),radial-gradient(circle at 70% 80%,rgba(16,185,129,0.02) 0%,transparent 40%);pointer-events:none;z-index:0;animation:bgShift 30s ease-in-out infinite alternate}
@keyframes bgShift{0%{transform:translate(0,0)}100%{transform:translate(2%,1%)}}

/* Scrollbar */
::-webkit-scrollbar{width:5px}
::-webkit-scrollbar-track{background:transparent}
::-webkit-scrollbar-thumb{background:var(--scrollbar-thumb);border-radius:10px}
::-webkit-scrollbar-thumb:hover{background:var(--scrollbar-hover)}

/* Layout */
.app{display:flex;min-height:100vh;position:relative;z-index:1}

/* Sidebar */
.sidebar{width:260px;background:var(--bg-sidebar);backdrop-filter:blur(20px) saturate(180%);-webkit-backdrop-filter:blur(20px) saturate(180%);border-right:1px solid var(--border);display:flex;flex-direction:column;position:fixed;top:0;left:0;bottom:0;z-index:100;transition:var(--transition-slow)}
.sidebar-header{padding:28px 24px;border-bottom:1px solid var(--border)}
.sidebar-logo{display:flex;align-items:center;gap:14px}
.sidebar-logo-icon{width:42px;height:42px;background:var(--gradient-1);border-radius:12px;display:flex;align-items:center;justify-content:center;font-size:20px;box-shadow:0 4px 16px var(--logo-glow);position:relative}
.sidebar-logo-icon::after{content:'';position:absolute;inset:-2px;border-radius:14px;background:var(--gradient-1);opacity:0.2;filter:blur(8px);z-index:-1}
.sidebar-logo h1{font-size:17px;font-weight:700;background:var(--gradient-1);-webkit-background-clip:text;-webkit-text-fill-color:transparent;letter-spacing:-0.3px}
.sidebar-logo p{font-size:11px;color:var(--text-muted);margin-top:3px;letter-spacing:0.3px}

.sidebar-nav{flex:1;padding:20px 14px;overflow-y:auto}
.nav-section{margin-bottom:28px}
.nav-section-title{font-size:10px;font-weight:700;text-transform:uppercase;letter-spacing:2px;color:var(--text-muted);padding:0 14px;margin-bottom:10px}
.nav-item{display:flex;align-items:center;gap:12px;padding:11px 14px;border-radius:var(--radius);color:var(--text-secondary);cursor:pointer;transition:var(--transition);font-size:13px;font-weight:500;position:relative;overflow:hidden;margin-bottom:2px}
.nav-item:hover{background:var(--accent-subtle);color:var(--text-primary)}
.nav-item.active{background:rgba(99,102,241,0.12);color:var(--accent-light);font-weight:600}
.nav-item.active::before{content:'';position:absolute;left:0;top:50%;transform:translateY(-50%);width:3px;height:55%;background:var(--gradient-1);border-radius:0 4px 4px 0}
.nav-icon{width:20px;text-align:center;font-size:15px;flex-shrink:0}

.sidebar-footer{padding:20px 24px;border-top:1px solid var(--border)}
.sidebar-status{display:flex;align-items:center;gap:10px;font-size:12px;color:var(--text-muted)}
.status-dot{width:8px;height:8px;border-radius:50%;background:var(--success);box-shadow:0 0 10px var(--success-glow),0 0 20px var(--success-glow);animation:pulse 2s ease-in-out infinite}
@keyframes pulse{0%,100%{opacity:1;transform:scale(1)}50%{opacity:0.6;transform:scale(0.9)}}

/* Main content */
.main{margin-left:260px;flex:1;padding:36px 40px;min-height:100vh}
.page{display:none;animation:fadeIn 0.4s cubic-bezier(0.4,0,0.2,1)}
.page.active{display:block}
@keyframes fadeIn{from{opacity:0;transform:translateY(12px)}to{opacity:1;transform:translateY(0)}}

/* Page header */
.page-header{margin-bottom:36px}
.page-header h2{font-size:26px;font-weight:800;margin-bottom:6px;letter-spacing:-0.5px}
.page-header p{color:var(--text-secondary);font-size:14px;letter-spacing:0.1px}

/* Cards */
.card{background:var(--bg-card);backdrop-filter:blur(10px);-webkit-backdrop-filter:blur(10px);border:1px solid var(--border);border-radius:var(--radius-lg);padding:28px;transition:var(--transition);box-shadow:var(--shadow-card)}
.card:hover{border-color:var(--border-light);box-shadow:var(--shadow-card),0 0 0 1px rgba(255,255,255,0.04)}
.card-header{display:flex;justify-content:space-between;align-items:center;margin-bottom:24px}
.card-title{font-size:16px;font-weight:700;letter-spacing:-0.2px}

/* Stat cards */
.stats-grid{display:grid;grid-template-columns:repeat(auto-fit,minmax(240px,1fr));gap:20px;margin-bottom:36px}
.stat-card{background:var(--bg-card);backdrop-filter:blur(10px);-webkit-backdrop-filter:blur(10px);border:1px solid var(--border);border-radius:var(--radius-lg);padding:28px;position:relative;overflow:hidden;transition:var(--transition)}
.stat-card:hover{transform:translateY(-3px);box-shadow:var(--shadow-md);border-color:var(--border-light)}
.stat-card::before{content:'';position:absolute;top:0;left:0;right:0;height:2px;opacity:0.8}
.stat-card::after{content:'';position:absolute;top:0;right:0;width:120px;height:120px;border-radius:50%;opacity:0.04;filter:blur(30px)}
.stat-card:nth-child(1)::before{background:var(--gradient-1)}
.stat-card:nth-child(1)::after{background:#6366f1}
.stat-card:nth-child(2)::before{background:var(--gradient-2)}
.stat-card:nth-child(2)::after{background:#10b981}
.stat-card:nth-child(3)::before{background:var(--gradient-3)}
.stat-card:nth-child(3)::after{background:#f43f5e}
.stat-card:nth-child(4)::before{background:var(--gradient-4)}
.stat-card:nth-child(4)::after{background:#f59e0b}
.stat-card-icon{width:48px;height:48px;border-radius:14px;display:flex;align-items:center;justify-content:center;font-size:22px;margin-bottom:20px}
.stat-card:nth-child(1) .stat-card-icon{background:rgba(99,102,241,0.1);color:var(--accent-light)}
.stat-card:nth-child(2) .stat-card-icon{background:rgba(16,185,129,0.1);color:var(--success)}
.stat-card:nth-child(3) .stat-card-icon{background:rgba(244,63,94,0.1);color:#fb7185}
.stat-card:nth-child(4) .stat-card-icon{background:rgba(245,158,11,0.1);color:var(--warning)}
.stat-card-value{font-size:32px;font-weight:800;margin-bottom:6px;font-variant-numeric:tabular-nums;letter-spacing:-1px}
.stat-card-label{font-size:13px;color:var(--text-secondary);font-weight:500}

/* Table */
.table-wrapper{overflow-x:auto;border-radius:var(--radius-lg);border:1px solid var(--border);background:var(--bg-table)}
table{width:100%;border-collapse:collapse}
th{padding:14px 20px;text-align:left;font-size:11px;font-weight:700;text-transform:uppercase;letter-spacing:1.2px;color:var(--text-muted);background:rgba(255,255,255,0.02);border-bottom:1px solid var(--border)}
td{padding:16px 20px;font-size:13px;border-bottom:1px solid rgba(255,255,255,0.03)}
tr:last-child td{border-bottom:none}
tr{transition:var(--transition)}
tr:hover td{background:var(--tr-hover)}

/* Badges */
.badge{display:inline-flex;align-items:center;padding:5px 12px;border-radius:8px;font-size:11px;font-weight:600;gap:5px;letter-spacing:0.2px}
.badge-success{background:rgba(16,185,129,0.1);color:var(--success);border:1px solid rgba(16,185,129,0.15)}
.badge-danger{background:rgba(239,68,68,0.1);color:var(--danger);border:1px solid rgba(239,68,68,0.15)}
.badge-warning{background:rgba(245,158,11,0.1);color:var(--warning);border:1px solid rgba(245,158,11,0.15)}
.badge-info{background:rgba(59,130,246,0.1);color:var(--info);border:1px solid rgba(59,130,246,0.15)}
.badge-purple{background:rgba(99,102,241,0.1);color:var(--accent-light);border:1px solid rgba(99,102,241,0.15)}
.badge-dot{width:6px;height:6px;border-radius:50%;display:inline-block}
.badge-dot.green{background:var(--success);box-shadow:0 0 6px var(--success-glow)}
.badge-dot.red{background:var(--danger);box-shadow:0 0 6px var(--danger-glow)}
.badge-dot.yellow{background:var(--warning)}

/* Buttons */
.btn{display:inline-flex;align-items:center;gap:7px;padding:9px 18px;border-radius:var(--radius);font-size:13px;font-weight:600;border:none;cursor:pointer;transition:var(--transition);font-family:inherit;letter-spacing:0.1px;position:relative;overflow:hidden}
.btn-primary{background:var(--gradient-1);color:#fff;box-shadow:0 2px 12px var(--accent-glow)}
.btn-primary:hover{transform:translateY(-1px);box-shadow:0 6px 24px var(--accent-glow);filter:brightness(1.1)}
.btn-primary:active{transform:translateY(0);filter:brightness(0.95)}
.btn-secondary{background:rgba(255,255,255,0.04);color:var(--text-primary);border:1px solid var(--border)}
.btn-secondary:hover{border-color:var(--border-accent);background:var(--accent-subtle)}
.btn-danger{background:rgba(239,68,68,0.1);color:var(--danger);border:1px solid rgba(239,68,68,0.15)}
.btn-danger:hover{background:rgba(239,68,68,0.2);border-color:rgba(239,68,68,0.3)}
.btn-sm{padding:7px 14px;font-size:12px;border-radius:8px}
.btn-icon{width:34px;height:34px;padding:0;display:flex;align-items:center;justify-content:center;border-radius:10px}

/* Forms */
.form-group{margin-bottom:22px}
.form-label{display:block;font-size:12px;font-weight:600;color:var(--text-secondary);margin-bottom:8px;text-transform:uppercase;letter-spacing:0.8px}
.form-input,.form-select,.form-textarea{width:100%;padding:11px 16px;background:rgba(255,255,255,0.03);border:1px solid var(--border);border-radius:var(--radius);color:var(--text-primary);font-size:13px;font-family:inherit;transition:var(--transition);outline:none}
.form-input:focus,.form-select:focus,.form-textarea:focus{border-color:var(--accent);box-shadow:0 0 0 3px var(--accent-glow),0 0 20px rgba(99,102,241,0.08)}
.form-input::placeholder{color:var(--text-muted)}
.form-textarea{min-height:80px;resize:vertical}
.form-select{appearance:none;background-image:url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='12' viewBox='0 0 12 12'%3E%3Cpath fill='currentColor' d='M6 8L1 3h10z'/%3E%3C/svg%3E");background-repeat:no-repeat;background-position:right 14px center;padding-right:36px}
.form-row{display:grid;grid-template-columns:1fr 1fr;gap:18px}
.form-hint{font-size:11px;color:var(--text-muted);margin-top:6px}

/* Modal */
.modal-overlay{display:none;position:fixed;top:0;left:0;right:0;bottom:0;background:rgba(0,0,0,0.7);backdrop-filter:blur(8px);-webkit-backdrop-filter:blur(8px);z-index:200;align-items:center;justify-content:center}
.modal-overlay.active{display:flex}
.modal{background:var(--bg-modal);border:1px solid var(--border-light);border-radius:var(--radius-xl);width:90%;max-width:580px;max-height:85vh;overflow-y:auto;box-shadow:var(--shadow-lg),var(--shadow-glow);animation:modalIn 0.35s cubic-bezier(0.34,1.56,0.64,1)}
@keyframes modalIn{from{opacity:0;transform:scale(0.92) translateY(16px)}to{opacity:1;transform:scale(1) translateY(0)}}
.modal-header{padding:28px 28px 0;display:flex;justify-content:space-between;align-items:center}
.modal-header h3{font-size:19px;font-weight:700;letter-spacing:-0.3px}
.modal-close{width:34px;height:34px;display:flex;align-items:center;justify-content:center;border-radius:10px;border:none;background:rgba(255,255,255,0.04);color:var(--text-muted);cursor:pointer;font-size:16px;transition:var(--transition)}
.modal-close:hover{background:rgba(239,68,68,0.12);color:var(--danger)}
.modal-body{padding:28px}
.modal-footer{padding:0 28px 28px;display:flex;justify-content:flex-end;gap:10px}

/* Toast */
.toast-container{position:fixed;top:24px;right:24px;z-index:300;display:flex;flex-direction:column;gap:10px}
.toast{display:flex;align-items:center;gap:12px;padding:16px 22px;border-radius:var(--radius);background:var(--bg-toast);backdrop-filter:blur(20px);-webkit-backdrop-filter:blur(20px);border:1px solid var(--border);box-shadow:var(--shadow-lg);font-size:13px;animation:toastIn 0.4s cubic-bezier(0.34,1.56,0.64,1);min-width:300px;font-weight:500}
.toast.success{border-color:rgba(16,185,129,0.25);box-shadow:var(--shadow-lg),0 0 20px rgba(16,185,129,0.08)}
.toast.error{border-color:rgba(239,68,68,0.25);box-shadow:var(--shadow-lg),0 0 20px rgba(239,68,68,0.08)}
@keyframes toastIn{from{opacity:0;transform:translateX(100%) scale(0.9)}to{opacity:1;transform:translateX(0) scale(1)}}

/* Login page */
.login-container{display:flex;align-items:center;justify-content:center;min-height:100vh;position:relative;z-index:1}
.login-container::before{content:'';position:fixed;top:50%;left:50%;width:600px;height:600px;transform:translate(-50%,-50%);background:radial-gradient(circle,var(--login-glow) 0%,transparent 70%);pointer-events:none;animation:loginGlow 8s ease-in-out infinite alternate}
@keyframes loginGlow{0%{opacity:0.5;transform:translate(-50%,-50%) scale(1)}100%{opacity:1;transform:translate(-50%,-50%) scale(1.2)}}
.login-card{background:var(--bg-login);backdrop-filter:blur(30px) saturate(180%);-webkit-backdrop-filter:blur(30px) saturate(180%);border:1px solid var(--border-light);border-radius:var(--radius-xl);padding:52px;width:100%;max-width:440px;box-shadow:var(--shadow-lg),var(--shadow-glow);position:relative}
.login-card::before{content:'';position:absolute;top:-1px;left:20%;right:20%;height:1px;background:linear-gradient(90deg,transparent,rgba(99,102,241,0.4),transparent)}
.login-card h2{font-size:26px;font-weight:800;text-align:center;margin-bottom:8px;background:var(--gradient-1);-webkit-background-clip:text;-webkit-text-fill-color:transparent;letter-spacing:-0.5px}
.login-card p{text-align:center;color:var(--text-secondary);margin-bottom:36px;font-size:14px}
.login-card .form-input{padding:13px 18px;background:rgba(255,255,255,0.03);font-size:14px}
.login-card .btn-primary{width:100%;padding:13px;font-size:14px;justify-content:center;margin-top:12px}

/* Dashboard chart area placeholder */
.chart-area{background:var(--bg-card);border:1px solid var(--border);border-radius:var(--radius-lg);padding:28px;min-height:200px}
.chart-placeholder{display:flex;flex-direction:column;align-items:center;justify-content:center;min-height:180px;color:var(--text-muted);gap:10px}

/* Provider type tags */
.provider-tag{display:inline-flex;align-items:center;gap:6px;padding:4px 12px;border-radius:8px;font-size:11px;font-weight:600;letter-spacing:0.2px}
.provider-tag.openai{background:var(--tag-openai-bg);color:var(--tag-openai-c);border:1px solid var(--tag-openai-b)}
.provider-tag.claude{background:var(--tag-claude-bg);color:var(--tag-claude-c);border:1px solid var(--tag-claude-b)}
.provider-tag.qwen{background:var(--tag-qwen-bg);color:var(--tag-qwen-c);border:1px solid var(--tag-qwen-b)}
.provider-tag.default{background:var(--bg-surface);color:var(--text-secondary);border:1px solid var(--border)}

/* Empty state */
.empty-state{text-align:center;padding:72px 20px;color:var(--text-muted)}
.empty-state-icon{font-size:52px;margin-bottom:20px;opacity:0.4;filter:grayscale(20%)}
.empty-state h3{font-size:16px;color:var(--text-secondary);margin-bottom:10px;font-weight:600}
.empty-state p{font-size:13px;max-width:340px;margin:0 auto;line-height:1.7}

/* Copy button */
.copy-btn{display:inline-flex;align-items:center;gap:5px;padding:5px 10px;border-radius:8px;border:1px solid var(--border);background:rgba(255,255,255,0.03);color:var(--text-secondary);cursor:pointer;font-size:11px;font-family:inherit;transition:var(--transition)}
.copy-btn:hover{border-color:var(--border-accent);color:var(--accent-light);background:var(--accent-subtle)}

/* Key display */
.key-display{background:rgba(255,255,255,0.02);border:1px solid var(--border);border-radius:var(--radius);padding:18px;font-family:'Courier New',monospace;font-size:13px;word-break:break-all;color:var(--success);margin:14px 0;position:relative}
.key-warning{display:flex;align-items:flex-start;gap:10px;padding:14px;background:rgba(245,158,11,0.06);border:1px solid rgba(245,158,11,0.12);border-radius:var(--radius);font-size:12px;color:var(--warning);margin-top:14px;line-height:1.6}

/* Data table (used in users) */
.data-table{width:100%;border-collapse:collapse}
.data-table thead tr{border-bottom:1px solid var(--border)}
.data-table th{padding:14px 20px;text-align:left;font-size:11px;font-weight:700;text-transform:uppercase;letter-spacing:1.2px;color:var(--text-muted);background:rgba(255,255,255,0.02)}
.data-table td{padding:16px 20px;font-size:13px;border-bottom:1px solid rgba(255,255,255,0.03)}
.data-table tr:hover td{background:var(--tr-hover)}
.status-badge{display:inline-flex;align-items:center;padding:4px 10px;border-radius:8px;font-size:11px;font-weight:600}
.status-badge.success{background:rgba(16,185,129,0.1);color:var(--success)}
.action-buttons{display:flex;gap:6px}

/* Responsive */
@media(max-width:768px){
  .sidebar{width:0;overflow:hidden}
  .main{margin-left:0;padding:24px 16px}
  .stats-grid{grid-template-columns:1fr}
  .form-row{grid-template-columns:1fr}
}
</style>
</head>
<body>

<!-- Toast container -->
<div class="toast-container" id="toastContainer"></div>

<!-- Login view -->
<div id="loginView" class="login-container">
  <div class="login-card">
    <div style="text-align:center;margin-bottom:32px">
      <div class="sidebar-logo-icon" style="margin:0 auto 20px;width:64px;height:64px;font-size:30px;border-radius:18px;box-shadow:0 8px 32px rgba(99,102,241,0.3)">⚡</div>
    </div>
    <h2>AI Gateway</h2>
    <p>统一 AI 模型管理平台</p>
    <div class="form-group">
      <label class="form-label">用户名</label>
      <input type="text" class="form-input" id="loginUser" value="admin" placeholder="请输入用户名">
    </div>
    <div class="form-group">
      <label class="form-label">密码</label>
      <input type="password" class="form-input" id="loginPass" placeholder="请输入密码" onkeypress="if(event.key==='Enter')doLogin()">
    </div>
    <button class="btn btn-primary" onclick="doLogin()">登录管理后台</button>
    <div style="text-align:center;margin-top:24px">
      <span style="font-size:11px;color:var(--text-muted);letter-spacing:0.3px">默认密码为 JWT Secret 前 8 位字符</span>
    </div>
  </div>
</div>

<!-- App view -->
<div id="appView" class="app" style="display:none">
  <!-- Sidebar -->
  <aside class="sidebar">
    <div class="sidebar-header">
      <div class="sidebar-logo">
        <div class="sidebar-logo-icon">⚡</div>
        <div>
          <h1>AI Gateway</h1>
          <p>Unified Model Platform</p>
        </div>
      </div>
    </div>
    <nav class="sidebar-nav">
      <div class="nav-section">
        <div class="nav-section-title">概览</div>
        <div class="nav-item active" onclick="showPage('dashboard')">
          <span class="nav-icon">📊</span>仪表盘
        </div>
      </div>
      <div class="nav-section">
        <div class="nav-section-title">管理</div>
        <div class="nav-item" onclick="showPage('providers')">
          <span class="nav-icon">🏢</span>厂商管理
        </div>
        <div class="nav-item" onclick="showPage('models')">
          <span class="nav-icon">🤖</span>模型管理
        </div>
        <div class="nav-item" onclick="showPage('channels')">
          <span class="nav-icon">🔗</span>通道健康
        </div>
      </div>
      <div class="nav-section">
        <div class="nav-section-title">用户</div>
        <div class="nav-item" onclick="showPage('users')">
          <span class="nav-icon">👥</span>用户管理
        </div>
        <div class="nav-item" onclick="showPage('apikeys')">
          <span class="nav-icon">🔑</span>API Key 管理
        </div>
      </div>
    </nav>
    <div class="sidebar-footer">
      <div class="sidebar-status">
        <span class="status-dot"></span>
        <span>System Online</span>
      </div>
    </div>
  </aside>

  <!-- Main content -->
  <main class="main">
    <!-- Dashboard -->
    <div class="page active" id="page-dashboard">
      <div class="page-header">
        <h2>仪表盘</h2>
        <p>AI Gateway 运行状态概览</p>
      </div>
      <div class="stats-grid" id="dashboardStats">
        <div class="stat-card">
          <div class="stat-card-icon">🏢</div>
          <div class="stat-card-value" id="stat-providers">0</div>
          <div class="stat-card-label">AI 厂商</div>
        </div>
        <div class="stat-card">
          <div class="stat-card-icon">🤖</div>
          <div class="stat-card-value" id="stat-models">0</div>
          <div class="stat-card-label">可用模型</div>
        </div>
        <div class="stat-card">
          <div class="stat-card-icon">🔑</div>
          <div class="stat-card-value" id="stat-keys">0</div>
          <div class="stat-card-label">API Keys</div>
        </div>
        <div class="stat-card">
          <div class="stat-card-icon">🔗</div>
          <div class="stat-card-value" id="stat-channels">0</div>
          <div class="stat-card-label">活跃通道</div>
        </div>
      </div>
      <div class="card">
        <div class="card-header">
          <span class="card-title">通道健康状态</span>
          <button class="btn btn-secondary btn-sm" onclick="loadDashboard()">🔄 刷新</button>
        </div>
        <div id="dashboardChannels">
          <div class="empty-state">
            <div class="empty-state-icon">📡</div>
            <h3>暂无通道数据</h3>
            <p>请先添加厂商和模型以创建通道</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Providers -->
    <div class="page" id="page-providers">
      <div class="page-header" style="display:flex;justify-content:space-between;align-items:flex-start">
        <div>
          <h2>厂商管理</h2>
          <p>管理 AI 模型提供商及其 API 配置</p>
        </div>
        <button class="btn btn-primary" onclick="showProviderModal()">✚ 添加厂商</button>
      </div>
      <div class="card">
        <div id="providersList">
          <div class="empty-state">
            <div class="empty-state-icon">🏢</div>
            <h3>暂无厂商</h3>
            <p>点击"添加厂商"按钮来配置您的第一个 AI 提供商</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Models -->
    <div class="page" id="page-models">
      <div class="page-header" style="display:flex;justify-content:space-between;align-items:flex-start">
        <div>
          <h2>模型管理</h2>
          <p>管理各厂商下的可用模型和通道配置</p>
        </div>
        <button class="btn btn-primary" onclick="showModelModal()">✚ 添加模型</button>
      </div>
      <div class="card">
        <div id="modelsList">
          <div class="empty-state">
            <div class="empty-state-icon">🤖</div>
            <h3>暂无模型</h3>
            <p>请先添加厂商后再配置模型</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Channels -->
    <div class="page" id="page-channels">
      <div class="page-header">
        <h2>通道健康</h2>
        <p>实时监控所有通道的健康状态和性能指标</p>
      </div>
      <div class="card">
        <div class="card-header">
          <span class="card-title">通道状态</span>
          <button class="btn btn-secondary btn-sm" onclick="loadChannels()">🔄 刷新</button>
        </div>
        <div id="channelsList">
          <div class="empty-state">
            <div class="empty-state-icon">🔗</div>
            <h3>暂无通道</h3>
            <p>通道会在添加模型后自动创建</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Users -->
    <div class="page" id="page-users">
      <div class="page-header" style="display:flex;justify-content:space-between;align-items:flex-start">
        <div>
          <h2>用户管理</h2>
          <p>管理系统注册用户</p>
        </div>
        <div style="display:flex;gap:12px;align-items:center">
          <input type="text" class="form-input" id="userSearchInput" placeholder="搜索用户名/邮箱" style="width:220px" onkeydown="if(event.key==='Enter')loadUsers()">
          <button class="btn btn-secondary" onclick="loadUsers()">🔍 搜索</button>
          <button class="btn btn-primary" onclick="showUserModal()">✚ 创建用户</button>
        </div>
      </div>
      <div class="card">
        <div id="usersList">
          <div class="empty-state">
            <div class="empty-state-icon">👥</div>
            <h3>暂无用户</h3>
            <p>系统中还没有注册用户</p>
          </div>
        </div>
        <div id="usersPagination" style="display:flex;justify-content:space-between;align-items:center;margin-top:16px"></div>
      </div>
    </div>

    <!-- API Keys -->
    <div class="page" id="page-apikeys">
      <div class="page-header" style="display:flex;justify-content:space-between;align-items:flex-start">
        <div>
          <h2>API Key 管理</h2>
          <p>创建和管理用户 API Key</p>
        </div>
        <button class="btn btn-primary" onclick="showKeyModal()">✚ 创建 Key</button>
      </div>
      <div class="card">
        <div id="keysList">
          <div class="empty-state">
            <div class="empty-state-icon">🔑</div>
            <h3>暂无 API Key</h3>
            <p>点击"创建 Key"按钮来生成新的 API Key</p>
          </div>
        </div>
      </div>
    </div>
  </main>
</div>

<!-- Provider Modal -->
<div class="modal-overlay" id="providerModal">
  <div class="modal">
    <div class="modal-header">
      <h3 id="providerModalTitle">添加厂商</h3>
      <button class="modal-close" onclick="closeModal('providerModal')">✕</button>
    </div>
    <div class="modal-body">
      <input type="hidden" id="providerEditId">
      <div class="form-group">
        <label class="form-label">厂商名称</label>
        <input type="text" class="form-input" id="providerName" placeholder="例如：OpenAI">
      </div>
      <div class="form-row">
        <div class="form-group">
          <label class="form-label">厂商类型</label>
          <select class="form-select" id="providerType">
            <option value="openai">OpenAI</option>
            <option value="claude">Claude (Anthropic)</option>
            <option value="qwen">通义千问 (Qwen)</option>
          </select>
        </div>
        <div class="form-group">
          <label class="form-label">Organization ID</label>
          <input type="text" class="form-input" id="providerOrgId" placeholder="可选">
        </div>
      </div>
      <div class="form-group">
        <label class="form-label">Base URL</label>
        <input type="text" class="form-input" id="providerBaseUrl" placeholder="https://api.openai.com">
      </div>
      <div class="form-group">
        <label class="form-label">API Key</label>
        <input type="password" class="form-input" id="providerApiKey" placeholder="sk-xxxxxxxxxxxx">
        <div class="form-hint">密钥将使用 AES-256-GCM 加密存储，保存后无法查看</div>
      </div>
    </div>
    <div class="modal-footer">
      <button class="btn btn-secondary" onclick="closeModal('providerModal')">取消</button>
      <button class="btn btn-primary" onclick="saveProvider()">保存</button>
    </div>
  </div>
</div>

<!-- Model Modal -->
<div class="modal-overlay" id="modelModal">
  <div class="modal">
    <div class="modal-header">
      <h3>添加模型</h3>
      <button class="modal-close" onclick="closeModal('modelModal')">✕</button>
    </div>
    <div class="modal-body">
      <div class="form-group">
        <label class="form-label">所属厂商</label>
        <select class="form-select" id="modelProviderSelect"></select>
      </div>
      <div class="form-row">
        <div class="form-group">
          <label class="form-label">对外模型名</label>
          <input type="text" class="form-input" id="modelName" placeholder="例如：gpt-4o">
          <div class="form-hint">用户调用时使用的模型名</div>
        </div>
        <div class="form-group">
          <label class="form-label">上游模型 ID</label>
          <input type="text" class="form-input" id="modelId" placeholder="例如：gpt-4o-2024-08-06">
          <div class="form-hint">厂商 API 的实际模型 ID</div>
        </div>
      </div>
      <div class="form-row">
        <div class="form-group">
          <label class="form-label">模型类型</label>
          <select class="form-select" id="modelType">
            <option value="chat">Chat (对话)</option>
            <option value="embedding">Embedding (向量)</option>
            <option value="image">Image (图像)</option>
          </select>
        </div>
        <div class="form-group">
          <label class="form-label">最大上下文 (tokens)</label>
          <input type="number" class="form-input" id="modelMaxTokens" placeholder="128000">
        </div>
      </div>
      <div class="form-row">
        <div class="form-group">
          <label class="form-label">负载权重</label>
          <input type="number" class="form-input" id="modelWeight" value="1" min="1">
          <div class="form-hint">值越大流量越多</div>
        </div>
        <div class="form-group">
          <label class="form-label">优先级</label>
          <input type="number" class="form-input" id="modelPriority" value="0" min="0">
          <div class="form-hint">0=主通道，数值越大优先级越低</div>
        </div>
      </div>
      <div class="form-row">
        <div class="form-group">
          <label class="form-label">输入价格 ($/1K tokens)</label>
          <input type="number" class="form-input" id="modelInputPrice" step="0.0001" placeholder="0.0025">
        </div>
        <div class="form-group">
          <label class="form-label">输出价格 ($/1K tokens)</label>
          <input type="number" class="form-input" id="modelOutputPrice" step="0.0001" placeholder="0.01">
        </div>
      </div>
    </div>
    <div class="modal-footer">
      <button class="btn btn-secondary" onclick="closeModal('modelModal')">取消</button>
      <button class="btn btn-primary" onclick="saveModel()">保存</button>
    </div>
  </div>
</div>

<!-- API Key Modal -->
<div class="modal-overlay" id="keyModal">
  <div class="modal">
    <div class="modal-header">
      <h3>创建 API Key</h3>
      <button class="modal-close" onclick="closeModal('keyModal')">✕</button>
    </div>
    <div class="modal-body">
      <div class="form-group">
        <label class="form-label">Key 名称</label>
        <input type="text" class="form-input" id="keyName" placeholder="例如：开发测试用">
      </div>
      <div class="form-row">
        <div class="form-group">
          <label class="form-label">每分钟限流</label>
          <input type="number" class="form-input" id="keyRateLimit" value="60" min="1">
        </div>
        <div class="form-group">
          <label class="form-label">Token 配额</label>
          <input type="number" class="form-input" id="keyQuota" value="0">
          <div class="form-hint">0 表示无限制</div>
        </div>
      </div>
      <div class="form-group">
        <label class="form-label">过期时间</label>
        <input type="datetime-local" class="form-input" id="keyExpires">
        <div class="form-hint">留空表示永不过期</div>
      </div>
    </div>
    <div class="modal-footer">
      <button class="btn btn-secondary" onclick="closeModal('keyModal')">取消</button>
      <button class="btn btn-primary" onclick="saveKey()">生成 Key</button>
    </div>
  </div>
</div>

<!-- Key Result Modal -->
<div class="modal-overlay" id="keyResultModal">
  <div class="modal">
    <div class="modal-header">
      <h3>🎉 API Key 创建成功</h3>
      <button class="modal-close" onclick="closeModal('keyResultModal')">✕</button>
    </div>
    <div class="modal-body">
      <p style="margin-bottom:12px;color:var(--text-secondary)">请复制并妥善保存您的 API Key：</p>
      <div class="key-display" id="newKeyDisplay"></div>
      <button class="btn btn-secondary btn-sm" onclick="copyKey()" style="margin-top:8px">📋 复制 Key</button>
      <div class="key-warning">
        <span>⚠️</span>
        <span>此 Key 仅展示这一次，关闭后将无法再次查看完整 Key。请务必先复制保存。</span>
      </div>
    </div>
    <div class="modal-footer">
      <button class="btn btn-primary" onclick="closeModal('keyResultModal')">我已保存</button>
    </div>
  </div>
</div>

<!-- User Modal -->
<div class="modal-overlay" id="userModal">
  <div class="modal">
    <div class="modal-header">
      <h3 id="userModalTitle">创建用户</h3>
      <button class="modal-close" onclick="closeModal('userModal')">✕</button>
    </div>
    <div class="modal-body">
      <input type="hidden" id="userEditId">
      <div class="form-group">
        <label class="form-label">用户名</label>
        <input type="text" class="form-input" id="userUsername" placeholder="至少3个字符">
      </div>
      <div class="form-group">
        <label class="form-label">邮箱</label>
        <input type="email" class="form-input" id="userEmail" placeholder="user@example.com">
      </div>
      <div class="form-group">
        <label class="form-label" id="userPasswordLabel">密码</label>
        <input type="password" class="form-input" id="userPassword" placeholder="至少6个字符">
        <div class="form-hint" id="userPasswordHint"></div>
      </div>
      <div class="form-row">
        <div class="form-group">
          <label class="form-label">角色</label>
          <select class="form-select" id="userRole">
            <option value="user">普通用户</option>
            <option value="admin">管理员</option>
          </select>
        </div>
        <div class="form-group">
          <label class="form-label">状态</label>
          <select class="form-select" id="userEnabled">
            <option value="true">启用</option>
            <option value="false">禁用</option>
          </select>
        </div>
      </div>
    </div>
    <div class="modal-footer">
      <button class="btn btn-secondary" onclick="closeModal('userModal')">取消</button>
      <button class="btn btn-primary" onclick="saveUser()">保存</button>
    </div>
  </div>
</div>

<script>
let adminToken = '';
const API = '';

// ── Theme ───────────────────────
function setTheme(name) {
  document.documentElement.setAttribute('data-theme', name);
  localStorage.setItem('ui-theme', name);
  document.querySelectorAll('.theme-option').forEach(d => {
    d.classList.toggle('active', d.getAttribute('data-theme') === name);
  });
}
function toggleThemePanel() {
  const fab = document.getElementById('themeFab');
  const panel = document.getElementById('themePanel');
  const isOpen = panel.classList.toggle('open');
  fab.classList.toggle('open', isOpen);
}
// Close panel on outside click
document.addEventListener('click', function(e) {
  const panel = document.getElementById('themePanel');
  const fab = document.getElementById('themeFab');
  if (panel && fab && !panel.contains(e.target) && !fab.contains(e.target)) {
    panel.classList.remove('open');
    fab.classList.remove('open');
  }
});
(function initTheme() {
  const saved = localStorage.getItem('ui-theme') || 'midnight';
  setTheme(saved);
})();

// ── Auth ────────────────────────
async function doLogin() {
  const user = document.getElementById('loginUser').value;
  const pass = document.getElementById('loginPass').value;
  try {
    const r = await fetch(API + '/admin/login', {
      method: 'POST',
      headers: {'Content-Type': 'application/json'},
      body: JSON.stringify({username: user, password: pass})
    });
    const d = await r.json();
    if (d.code === 0 && d.data && d.data.token) {
      adminToken = d.data.token;
      localStorage.setItem('adminToken', adminToken);
      document.getElementById('loginView').style.display = 'none';
      document.getElementById('appView').style.display = 'flex';
      loadDashboard();
      toast('登录成功', 'success');
    } else {
      toast(d.error?.message || '登录失败', 'error');
    }
  } catch(e) {
    toast('连接失败: ' + e.message, 'error');
  }
}

function authHeaders() {
  return {'Authorization': 'Bearer ' + adminToken, 'Content-Type': 'application/json'};
}

// ── Init ────────────────────────
(function init() {
  const t = localStorage.getItem('adminToken');
  if (t) {
    adminToken = t;
    document.getElementById('loginView').style.display = 'none';
    document.getElementById('appView').style.display = 'flex';
    loadDashboard();
  }
})();

// ── Navigation ──────────────────
function showPage(name) {
  document.querySelectorAll('.page').forEach(p => p.classList.remove('active'));
  document.querySelectorAll('.nav-item').forEach(n => n.classList.remove('active'));
  document.getElementById('page-' + name).classList.add('active');
  event.currentTarget.classList.add('active');
  switch(name) {
    case 'dashboard': loadDashboard(); break;
    case 'providers': loadProviders(); break;
    case 'models': loadModels(); break;
    case 'channels': loadChannels(); break;
    case 'users': loadUsers(); break;
    case 'apikeys': loadKeys(); break;
  }
}

// ── Dashboard ───────────────────
async function loadDashboard() {
  try {
    const [pr, ch] = await Promise.all([
      fetch(API + '/admin/providers', {headers: authHeaders()}).then(r=>r.json()),
      fetch(API + '/admin/channels/health', {headers: authHeaders()}).then(r=>r.json())
    ]);

    const providers = pr.data || [];
    const channels = ch.data || [];

    document.getElementById('stat-providers').textContent = providers.length;

    // Count unique model names
    const modelNames = new Set();
    channels.forEach(c => modelNames.add(c.model_name));
    document.getElementById('stat-models').textContent = modelNames.size;

    document.getElementById('stat-channels').textContent = channels.length;

    // Load keys count
    try {
      const kr = await fetch(API + '/admin/api/keys', {headers: authHeaders()}).then(r=>r.json());
      document.getElementById('stat-keys').textContent = (kr.data || []).length;
    } catch(e) {
      document.getElementById('stat-keys').textContent = '-';
    }

    // Channel health table
    if (channels.length === 0) {
      document.getElementById('dashboardChannels').innerHTML =
        '<div class="empty-state"><div class="empty-state-icon">📡</div><h3>暂无通道数据</h3><p>请先添加厂商和模型</p></div>';
    } else {
      let html = '<div class="table-wrapper"><table><thead><tr><th>通道 ID</th><th>模型</th><th>厂商</th><th>权重</th><th>优先级</th><th>状态</th><th>连续失败</th></tr></thead><tbody>';
      channels.forEach(c => {
        const statusClass = c.status === 'closed' ? 'success' : c.status === 'open' ? 'danger' : 'warning';
        const statusText = c.status === 'closed' ? '健康' : c.status === 'open' ? '熔断' : '半开';
        html += '<tr>';
        html += '<td><span style="font-weight:600">#' + c.channel_id + '</span></td>';
        html += '<td><span class="badge badge-purple">' + esc(c.model_name || '') + '</span></td>';
        html += '<td>' + esc(c.provider_name || '-') + '</td>';
        html += '<td>' + c.weight + '</td>';
        html += '<td>' + c.priority + '</td>';
        html += '<td><span class="badge badge-' + statusClass + '"><span class="badge-dot ' + (statusClass==='success'?'green':statusClass==='danger'?'red':'yellow') + '"></span>' + statusText + '</span></td>';
        html += '<td>' + (c.consecutive_failures || 0) + '</td>';
        html += '</tr>';
      });
      html += '</tbody></table></div>';
      document.getElementById('dashboardChannels').innerHTML = html;
    }
  } catch(e) {
    console.error(e);
  }
}

// ── Providers ───────────────────
async function loadProviders() {
  try {
    const r = await fetch(API + '/admin/providers', {headers: authHeaders()});
    const d = await r.json();
    const providers = d.data || [];
    if (providers.length === 0) {
      document.getElementById('providersList').innerHTML =
        '<div class="empty-state"><div class="empty-state-icon">🏢</div><h3>暂无厂商</h3><p>点击"添加厂商"按钮来配置</p></div>';
      return;
    }
    let html = '<div class="table-wrapper"><table><thead><tr><th>ID</th><th>名称</th><th>类型</th><th>Base URL</th><th>状态</th><th>操作</th></tr></thead><tbody>';
    providers.forEach(p => {
      const typeTag = providerTag(p.type);
      html += '<tr>';
      html += '<td><span style="font-weight:600">#' + p.id + '</span></td>';
      html += '<td style="font-weight:600">' + esc(p.name) + '</td>';
      html += '<td>' + typeTag + '</td>';
      html += '<td style="font-size:12px;color:var(--text-secondary)">' + esc(p.base_url) + '</td>';
      html += '<td>' + (p.enabled ? '<span class="badge badge-success"><span class="badge-dot green"></span>启用</span>' : '<span class="badge badge-danger"><span class="badge-dot red"></span>禁用</span>') + '</td>';
      html += '<td><button class="btn btn-secondary btn-sm" onclick="editProvider('+p.id+')" style="margin-right:4px">编辑</button>';
      html += '<button class="btn btn-danger btn-sm" onclick="deleteProvider('+p.id+')">删除</button></td>';
      html += '</tr>';
    });
    html += '</tbody></table></div>';
    document.getElementById('providersList').innerHTML = html;
  } catch(e) { toast('加载厂商失败', 'error'); }
}

let cachedProviders = [];
function showProviderModal(editing) {
  if (!editing) {
    document.getElementById('providerEditId').value = '';
    document.getElementById('providerName').value = '';
    document.getElementById('providerType').value = 'openai';
    document.getElementById('providerBaseUrl').value = '';
    document.getElementById('providerApiKey').value = '';
    document.getElementById('providerOrgId').value = '';
    document.getElementById('providerModalTitle').textContent = '添加厂商';
  }
  openModal('providerModal');
}

async function editProvider(id) {
  try {
    const r = await fetch(API + '/admin/providers', {headers: authHeaders()});
    const d = await r.json();
    const p = (d.data || []).find(x => x.id === id);
    if (!p) return toast('厂商不存在', 'error');
    document.getElementById('providerEditId').value = p.id;
    document.getElementById('providerName').value = p.name;
    document.getElementById('providerType').value = p.type;
    document.getElementById('providerBaseUrl').value = p.base_url;
    document.getElementById('providerApiKey').value = '';
    document.getElementById('providerOrgId').value = p.org_id || '';
    document.getElementById('providerModalTitle').textContent = '编辑厂商';
    showProviderModal(true);
  } catch(e) { toast('加载失败', 'error'); }
}

async function saveProvider() {
  const editId = document.getElementById('providerEditId').value;
  const body = {
    name: document.getElementById('providerName').value,
    type: document.getElementById('providerType').value,
    base_url: document.getElementById('providerBaseUrl').value,
    org_id: document.getElementById('providerOrgId').value
  };
  const apiKey = document.getElementById('providerApiKey').value;

  if (!body.name || !body.base_url) {
    toast('请填写必填字段', 'error');
    return;
  }

  try {
    let r;
    if (editId) {
      const updateBody = {...body};
      if (apiKey) updateBody.api_key = apiKey;
      r = await fetch(API + '/admin/providers/' + editId, {
        method: 'PUT', headers: authHeaders(), body: JSON.stringify(updateBody)
      });
    } else {
      if (!apiKey) { toast('请输入 API Key', 'error'); return; }
      body.api_key = apiKey;
      r = await fetch(API + '/admin/providers', {
        method: 'POST', headers: authHeaders(), body: JSON.stringify(body)
      });
    }
    const d = await r.json();
    if (d.code === 0) {
      toast(editId ? '厂商已更新' : '厂商已添加', 'success');
      closeModal('providerModal');
      loadProviders();
    } else {
      toast(d.error?.message || '操作失败', 'error');
    }
  } catch(e) { toast('操作失败: ' + e.message, 'error'); }
}

async function deleteProvider(id) {
  if (!confirm('确定要删除此厂商吗？关联的模型也会被删除。')) return;
  try {
    const r = await fetch(API + '/admin/providers/' + id, {method: 'DELETE', headers: authHeaders()});
    const d = await r.json();
    if (d.code === 0) { toast('厂商已删除', 'success'); loadProviders(); }
    else toast(d.error?.message || '删除失败', 'error');
  } catch(e) { toast('删除失败', 'error'); }
}

// ── Models ──────────────────────
async function loadModels() {
  try {
    const pr = await fetch(API + '/admin/providers', {headers: authHeaders()});
    const pd = await pr.json();
    cachedProviders = pd.data || [];

    let allModels = [];
    for (const p of cachedProviders) {
      const mr = await fetch(API + '/admin/providers/' + p.id + '/models', {headers: authHeaders()});
      const md = await mr.json();
      (md.data || []).forEach(m => {
        m._provider_name = p.name;
        m._provider_type = p.type;
        allModels.push(m);
      });
    }

    if (allModels.length === 0) {
      document.getElementById('modelsList').innerHTML =
        '<div class="empty-state"><div class="empty-state-icon">🤖</div><h3>暂无模型</h3><p>请先添加厂商后再配置模型</p></div>';
      return;
    }

    let html = '<div class="table-wrapper"><table><thead><tr><th>ID</th><th>模型名</th><th>上游 ID</th><th>厂商</th><th>类型</th><th>权重</th><th>优先级</th><th>状态</th><th>操作</th></tr></thead><tbody>';
    allModels.forEach(m => {
      html += '<tr>';
      html += '<td><span style="font-weight:600">#' + m.id + '</span></td>';
      html += '<td><span class="badge badge-purple">' + esc(m.model_name) + '</span></td>';
      html += '<td style="font-size:12px;color:var(--text-secondary)">' + esc(m.model_id) + '</td>';
      html += '<td>' + providerTag(m._provider_type) + ' ' + esc(m._provider_name) + '</td>';
      html += '<td><span class="badge badge-info">' + esc(m.model_type || 'chat') + '</span></td>';
      html += '<td style="font-weight:600">' + m.weight + '</td>';
      html += '<td>' + m.priority + '</td>';
      html += '<td>' + (m.enabled ? '<span class="badge badge-success"><span class="badge-dot green"></span>启用</span>' : '<span class="badge badge-danger"><span class="badge-dot red"></span>禁用</span>') + '</td>';
      html += '<td><button class="btn btn-danger btn-sm" onclick="deleteModel('+m.id+')">删除</button></td>';
      html += '</tr>';
    });
    html += '</tbody></table></div>';
    document.getElementById('modelsList').innerHTML = html;
  } catch(e) { toast('加载模型失败', 'error'); }
}

async function showModelModal() {
  try {
    const r = await fetch(API + '/admin/providers', {headers: authHeaders()});
    const d = await r.json();
    cachedProviders = d.data || [];
    const sel = document.getElementById('modelProviderSelect');
    sel.innerHTML = '';
    cachedProviders.forEach(p => {
      sel.innerHTML += '<option value="' + p.id + '">' + esc(p.name) + ' (' + p.type + ')</option>';
    });
    if (cachedProviders.length === 0) {
      toast('请先添加厂商', 'error');
      return;
    }
    document.getElementById('modelName').value = '';
    document.getElementById('modelId').value = '';
    document.getElementById('modelWeight').value = '1';
    document.getElementById('modelPriority').value = '0';
    document.getElementById('modelMaxTokens').value = '';
    document.getElementById('modelInputPrice').value = '';
    document.getElementById('modelOutputPrice').value = '';
    openModal('modelModal');
  } catch(e) { toast('加载厂商列表失败', 'error'); }
}

async function saveModel() {
  const providerId = document.getElementById('modelProviderSelect').value;
  const body = {
    model_name: document.getElementById('modelName').value,
    model_id: document.getElementById('modelId').value,
    model_type: document.getElementById('modelType').value,
    weight: parseInt(document.getElementById('modelWeight').value) || 1,
    priority: parseInt(document.getElementById('modelPriority').value) || 0
  };
  const maxTokens = document.getElementById('modelMaxTokens').value;
  if (maxTokens) body.max_context_tokens = parseInt(maxTokens);
  const inputPrice = document.getElementById('modelInputPrice').value;
  if (inputPrice) body.input_price = parseFloat(inputPrice);
  const outputPrice = document.getElementById('modelOutputPrice').value;
  if (outputPrice) body.output_price = parseFloat(outputPrice);

  if (!body.model_name || !body.model_id) {
    toast('请填写模型名称和上游 ID', 'error');
    return;
  }

  try {
    const r = await fetch(API + '/admin/providers/' + providerId + '/models', {
      method: 'POST', headers: authHeaders(), body: JSON.stringify(body)
    });
    const d = await r.json();
    if (d.code === 0) {
      toast('模型已添加', 'success');
      closeModal('modelModal');
      loadModels();
    } else {
      toast(d.error?.message || '添加失败', 'error');
    }
  } catch(e) { toast('添加失败: ' + e.message, 'error'); }
}

async function deleteModel(id) {
  if (!confirm('确定要删除此模型通道吗？')) return;
  try {
    const r = await fetch(API + '/admin/models/' + id, {method: 'DELETE', headers: authHeaders()});
    const d = await r.json();
    if (d.code === 0) { toast('模型已删除', 'success'); loadModels(); }
    else toast(d.error?.message || '删除失败', 'error');
  } catch(e) { toast('删除失败', 'error'); }
}

// ── Channels ────────────────────
async function loadChannels() {
  try {
    const r = await fetch(API + '/admin/channels/health', {headers: authHeaders()});
    const d = await r.json();
    const channels = d.data || [];
    if (channels.length === 0) {
      document.getElementById('channelsList').innerHTML =
        '<div class="empty-state"><div class="empty-state-icon">🔗</div><h3>暂无通道</h3><p>通道会在添加模型后自动创建</p></div>';
      return;
    }
    let html = '<div class="table-wrapper"><table><thead><tr><th>通道 ID</th><th>模型</th><th>厂商</th><th>权重</th><th>优先级</th><th>状态</th><th>连续失败</th><th>操作</th></tr></thead><tbody>';
    channels.forEach(c => {
      const statusClass = c.status === 'closed' ? 'success' : c.status === 'open' ? 'danger' : 'warning';
      const statusText = c.status === 'closed' ? '健康' : c.status === 'open' ? '熔断' : '半开';
      html += '<tr>';
      html += '<td><span style="font-weight:600">#' + c.channel_id + '</span></td>';
      html += '<td><span class="badge badge-purple">' + esc(c.model_name || '') + '</span></td>';
      html += '<td>' + esc(c.provider_name || '-') + '</td>';
      html += '<td>' + c.weight + '</td>';
      html += '<td>' + c.priority + '</td>';
      html += '<td><span class="badge badge-' + statusClass + '"><span class="badge-dot ' + (statusClass==='success'?'green':statusClass==='danger'?'red':'yellow') + '"></span>' + statusText + '</span></td>';
      html += '<td>' + (c.consecutive_failures || 0) + '</td>';
      html += '<td><button class="btn btn-secondary btn-sm" onclick="resetChannel('+c.channel_id+')">🔄 重置</button></td>';
      html += '</tr>';
    });
    html += '</tbody></table></div>';
    document.getElementById('channelsList').innerHTML = html;
  } catch(e) { toast('加载通道失败', 'error'); }
}

async function resetChannel(id) {
  try {
    await fetch(API + '/admin/channels/' + id + '/reset', {method: 'POST', headers: authHeaders()});
    toast('通道已重置', 'success');
    loadChannels();
  } catch(e) { toast('重置失败', 'error'); }
}

// ── API Keys ────────────────────
async function loadKeys() {
  try {
    const r = await fetch(API + '/admin/api/keys', {headers: authHeaders()});
    const d = await r.json();
    const keys = d.data || [];
    if (keys.length === 0) {
      document.getElementById('keysList').innerHTML =
        '<div class="empty-state"><div class="empty-state-icon">🔑</div><h3>暂无 API Key</h3><p>点击"创建 Key"按钮来生成</p></div>';
      return;
    }
    let html = '<div class="table-wrapper"><table><thead><tr><th>ID</th><th>名称</th><th>Key 前缀</th><th>限流</th><th>配额</th><th>状态</th><th>最后使用</th><th>操作</th></tr></thead><tbody>';
    keys.forEach(k => {
      const quotaText = k.quota_limit === 0 ? '无限' : (k.quota_used + ' / ' + k.quota_limit);
      html += '<tr>';
      html += '<td><span style="font-weight:600">#' + k.id + '</span></td>';
      html += '<td>' + esc(k.name || '-') + '</td>';
      html += '<td><code style="background:var(--bg-input);padding:2px 8px;border-radius:4px;font-size:12px">' + esc(k.key_prefix) + '****</code></td>';
      html += '<td>' + k.rate_limit + '/min</td>';
      html += '<td>' + quotaText + '</td>';
      html += '<td>' + (k.enabled ? '<span class="badge badge-success"><span class="badge-dot green"></span>启用</span>' : '<span class="badge badge-danger"><span class="badge-dot red"></span>禁用</span>') + '</td>';
      html += '<td style="font-size:12px;color:var(--text-secondary)">' + (k.last_used_at ? new Date(k.last_used_at).toLocaleString('zh-CN') : '从未使用') + '</td>';
      html += '<td><button class="btn btn-danger btn-sm" onclick="deleteKey('+k.id+')">删除</button></td>';
      html += '</tr>';
    });
    html += '</tbody></table></div>';
    document.getElementById('keysList').innerHTML = html;
  } catch(e) { toast('加载 Key 列表失败', 'error'); }
}

function showKeyModal() {
  document.getElementById('keyName').value = '';
  document.getElementById('keyRateLimit').value = '60';
  document.getElementById('keyQuota').value = '0';
  document.getElementById('keyExpires').value = '';
  openModal('keyModal');
}

let lastCreatedKey = '';
async function saveKey() {
  const body = {
    name: document.getElementById('keyName').value,
    rate_limit: parseInt(document.getElementById('keyRateLimit').value) || 60,
    quota_limit: parseInt(document.getElementById('keyQuota').value) || 0
  };
  const exp = document.getElementById('keyExpires').value;
  if (exp) body.expires_at = new Date(exp).toISOString();

  try {
    const r = await fetch(API + '/admin/api/keys', {
      method: 'POST', headers: authHeaders(), body: JSON.stringify(body)
    });
    const d = await r.json();
    if (d.code === 0 && d.data) {
      lastCreatedKey = d.data.key;
      closeModal('keyModal');
      document.getElementById('newKeyDisplay').textContent = lastCreatedKey;
      openModal('keyResultModal');
      loadKeys();
    } else {
      toast(d.error?.message || '创建失败', 'error');
    }
  } catch(e) { toast('创建失败: ' + e.message, 'error'); }
}

async function deleteKey(id) {
  if (!confirm('确定要删除此 API Key 吗？')) return;
  try {
    const r = await fetch(API + '/admin/api/keys/' + id, {method: 'DELETE', headers: authHeaders()});
    const d = await r.json();
    if (d.code === 0) { toast('Key 已删除', 'success'); loadKeys(); }
    else toast(d.error?.message || '删除失败', 'error');
  } catch(e) { toast('删除失败', 'error'); }
}

function copyKey() {
  navigator.clipboard.writeText(lastCreatedKey).then(() => toast('已复制到剪贴板', 'success'));
}

// ── Users ───────────────────────
let usersPage = 1;
const usersPageSize = 20;

async function loadUsers(page) {
  if (page) usersPage = page;
  const keyword = document.getElementById('userSearchInput').value || '';
  try {
    const r = await fetch(API + '/admin/users?page=' + usersPage + '&page_size=' + usersPageSize + '&keyword=' + encodeURIComponent(keyword), {headers: authHeaders()});
    const d = await r.json();
    const users = (d.data && d.data.users) || [];
    const total = (d.data && d.data.total) || 0;
    const el = document.getElementById('usersList');
    if (!users.length) {
      el.innerHTML = '<div class="empty-state"><div class="empty-state-icon">👥</div><h3>暂无用户</h3><p>没有找到匹配的用户</p></div>';
      document.getElementById('usersPagination').innerHTML = '';
      return;
    }
    let html = '<table class="data-table"><thead><tr><th>ID</th><th>用户名</th><th>邮箱</th><th>角色</th><th>状态</th><th>创建时间</th><th>操作</th></tr></thead><tbody>';
    users.forEach(u => {
      const roleBadge = u.role === 'admin' ? '<span style="color:var(--accent-light);font-weight:600">管理员</span>' : '<span style="color:var(--text-secondary)">用户</span>';
      const statusBadge = u.enabled ? '<span class="status-badge success">启用</span>' : '<span class="status-badge" style="background:var(--danger-glow);color:var(--danger)">禁用</span>';
      html += '<tr><td>' + u.id + '</td><td>' + esc(u.username) + '</td><td>' + esc(u.email) + '</td><td>' + roleBadge + '</td><td>' + statusBadge + '</td><td>' + esc(u.created_at) + '</td>';
      html += '<td><div class="action-buttons">';
      html += '<button class="btn btn-secondary btn-sm" onclick="editUser(' + u.id + ',\'' + esc(u.username) + '\',\'' + esc(u.email) + '\',\'' + u.role + '\',' + u.enabled + ')">编辑</button> ';
      html += '<button class="btn btn-sm" style="background:var(--danger-glow);color:var(--danger)" onclick="deleteUser(' + u.id + ',\'' + esc(u.username) + '\')">删除</button>';
      html += '</div></td></tr>';
    });
    html += '</tbody></table>';
    el.innerHTML = html;
    // Pagination
    const totalPages = Math.ceil(total / usersPageSize);
    let pagHtml = '<span style="color:var(--text-secondary);font-size:13px">共 ' + total + ' 条 · 第 ' + usersPage + '/' + totalPages + ' 页</span><div style="display:flex;gap:8px">';
    if (usersPage > 1) pagHtml += '<button class="btn btn-secondary btn-sm" onclick="loadUsers(' + (usersPage-1) + ')">上一页</button>';
    if (usersPage < totalPages) pagHtml += '<button class="btn btn-secondary btn-sm" onclick="loadUsers(' + (usersPage+1) + ')">下一页</button>';
    pagHtml += '</div>';
    document.getElementById('usersPagination').innerHTML = pagHtml;
  } catch(e) { toast('加载用户列表失败', 'error'); }
}

function showUserModal() {
  document.getElementById('userModalTitle').textContent = '创建用户';
  document.getElementById('userEditId').value = '';
  document.getElementById('userUsername').value = '';
  document.getElementById('userEmail').value = '';
  document.getElementById('userPassword').value = '';
  document.getElementById('userPasswordHint').textContent = '';
  document.getElementById('userPasswordLabel').textContent = '密码';
  document.getElementById('userRole').value = 'user';
  document.getElementById('userEnabled').value = 'true';
  openModal('userModal');
}

function editUser(id, username, email, role, enabled) {
  document.getElementById('userModalTitle').textContent = '编辑用户';
  document.getElementById('userEditId').value = id;
  document.getElementById('userUsername').value = username;
  document.getElementById('userEmail').value = email;
  document.getElementById('userPassword').value = '';
  document.getElementById('userPasswordHint').textContent = '留空表示不修改密码';
  document.getElementById('userPasswordLabel').textContent = '新密码（可选）';
  document.getElementById('userRole').value = role;
  document.getElementById('userEnabled').value = enabled ? 'true' : 'false';
  openModal('userModal');
}

async function saveUser() {
  const id = document.getElementById('userEditId').value;
  const username = document.getElementById('userUsername').value.trim();
  const email = document.getElementById('userEmail').value.trim();
  const password = document.getElementById('userPassword').value;
  const role = document.getElementById('userRole').value;
  const enabled = document.getElementById('userEnabled').value === 'true';

  if (!username || !email) { toast('请填写用户名和邮箱', 'error'); return; }

  try {
    let url, method, body;
    if (id) {
      url = API + '/admin/users/' + id;
      method = 'PUT';
      body = {username, email, role, enabled};
      if (password) body.password = password;
    } else {
      if (!password) { toast('请填写密码', 'error'); return; }
      url = API + '/admin/users';
      method = 'POST';
      body = {username, email, password, role};
    }
    const r = await fetch(url, {method, headers: {...authHeaders(), 'Content-Type':'application/json'}, body: JSON.stringify(body)});
    const d = await r.json();
    if (!r.ok) { toast(d.error?.message || '操作失败', 'error'); return; }
    toast(id ? '用户已更新' : '用户已创建', 'success');
    closeModal('userModal');
    loadUsers();
  } catch(e) { toast('操作失败: ' + e.message, 'error'); }
}

async function deleteUser(id, username) {
  if (!confirm('确定要删除用户 "' + username + '" 吗？此操作不可恢复。')) return;
  try {
    const r = await fetch(API + '/admin/users/' + id, {method:'DELETE', headers: authHeaders()});
    if (!r.ok) { const d = await r.json(); toast(d.error?.message || '删除失败', 'error'); return; }
    toast('用户已删除', 'success');
    loadUsers();
  } catch(e) { toast('删除失败', 'error'); }
}

// ── Utils ───────────────────────
function openModal(id) { document.getElementById(id).classList.add('active'); }
function closeModal(id) { document.getElementById(id).classList.remove('active'); }

function esc(s) {
  const d = document.createElement('div');
  d.textContent = s;
  return d.innerHTML;
}

function providerTag(type) {
  const cls = ['openai','claude','qwen'].includes(type) ? type : 'default';
  const labels = {openai:'OpenAI', claude:'Claude', qwen:'Qwen'};
  return '<span class="provider-tag ' + cls + '">' + (labels[type] || type) + '</span>';
}

function toast(msg, type) {
  const c = document.getElementById('toastContainer');
  const t = document.createElement('div');
  t.className = 'toast ' + (type || '');
  const icon = type === 'success' ? '✅' : type === 'error' ? '❌' : 'ℹ️';
  t.innerHTML = '<span>' + icon + '</span><span>' + esc(msg) + '</span>';
  c.appendChild(t);
  setTimeout(() => { t.style.opacity = '0'; t.style.transform = 'translateX(100%)'; setTimeout(() => t.remove(), 300); }, 3000);
}

// Close modals on overlay click
document.querySelectorAll('.modal-overlay').forEach(el => {
  el.addEventListener('click', e => { if (e.target === el) el.classList.remove('active'); });
});

// Auto-fill base URL based on type
document.getElementById('providerType').addEventListener('change', function() {
  const urls = {openai:'https://api.openai.com', claude:'https://api.anthropic.com', qwen:'https://dashscope.aliyuncs.com/compatible-mode'};
  const urlInput = document.getElementById('providerBaseUrl');
  if (!urlInput.value || Object.values(urls).includes(urlInput.value)) {
    urlInput.value = urls[this.value] || '';
  }
});
</script>

<!-- Theme Floating Button + Panel -->
<button class="theme-fab" id="themeFab" onclick="toggleThemePanel()" title="切换主题">🎨</button>
<div class="theme-panel" id="themePanel">
  <div class="theme-panel-title">选择主题</div>
  <div class="theme-option" data-theme="midnight" onclick="setTheme('midnight')">
    <div class="theme-option-colors">
      <div class="theme-option-dot" style="background:#6366f1"></div>
      <div class="theme-option-dot" style="background:#8b5cf6"></div>
    </div>
    <div class="theme-option-info">
      <div class="theme-option-name">Midnight</div>
      <div class="theme-option-desc">深邃暗黑 · 紫蓝</div>
    </div>
  </div>
  <div class="theme-option" data-theme="cyberpunk" onclick="setTheme('cyberpunk')">
    <div class="theme-option-colors">
      <div class="theme-option-dot" style="background:#00ffc8"></div>
      <div class="theme-option-dot" style="background:#00d4ff"></div>
    </div>
    <div class="theme-option-info">
      <div class="theme-option-name">Cyberpunk</div>
      <div class="theme-option-desc">赛博朋克 · 荧光青</div>
    </div>
  </div>
  <div class="theme-option" data-theme="aurora" onclick="setTheme('aurora')">
    <div class="theme-option-colors">
      <div class="theme-option-dot" style="background:#a855f7"></div>
      <div class="theme-option-dot" style="background:#ec4899"></div>
    </div>
    <div class="theme-option-info">
      <div class="theme-option-name">Aurora</div>
      <div class="theme-option-desc">极光梦境 · 紫粉</div>
    </div>
  </div>
  <div class="theme-option" data-theme="daylight" onclick="setTheme('daylight')">
    <div class="theme-option-colors">
      <div class="theme-option-dot" style="background:#f8f9fc;border:1px solid #ddd"></div>
      <div class="theme-option-dot" style="background:#6366f1"></div>
    </div>
    <div class="theme-option-info">
      <div class="theme-option-name">Daylight</div>
      <div class="theme-option-desc">明亮白昼 · 清爽</div>
    </div>
  </div>
  <div class="theme-option" data-theme="sunset" onclick="setTheme('sunset')">
    <div class="theme-option-colors">
      <div class="theme-option-dot" style="background:#f97316"></div>
      <div class="theme-option-dot" style="background:#ef4444"></div>
    </div>
    <div class="theme-option-info">
      <div class="theme-option-name">Sunset</div>
      <div class="theme-option-desc">日暮余晖 · 橙红</div>
    </div>
  </div>
</div>

</body>
</html>`
