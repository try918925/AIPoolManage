package main

const userPortalHTML = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>AI Gateway · 用户中心</title>
<link href="https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700;800&display=swap" rel="stylesheet">
<style>
*,*::before,*::after{margin:0;padding:0;box-sizing:border-box}
:root,
[data-theme="midnight"]{
  --bg-primary:#050509;
  --bg-secondary:#0c0c14;
  --bg-card:#111119;
  --bg-card-hover:#18182a;
  --bg-input:#0e0e1a;
  --bg-elevated:#1a1a3e;
  --border:rgba(255,255,255,0.06);
  --border-hover:rgba(255,255,255,0.1);
  --border-active:rgba(99,102,241,0.5);
  --text-primary:#f0f0f8;
  --text-secondary:#8b8ba8;
  --text-muted:#55556e;
  --accent:#6366f1;
  --accent-light:#818cf8;
  --accent-glow:rgba(99,102,241,0.25);
  --accent-soft:rgba(99,102,241,0.08);
  --accent-subtle:rgba(99,102,241,0.08);
  --success:#10b981;
  --success-soft:rgba(16,185,129,0.12);
  --warning:#f59e0b;
  --warning-soft:rgba(245,158,11,0.12);
  --danger:#ef4444;
  --danger-soft:rgba(239,68,68,0.12);
  --info:#3b82f6;
  --info-soft:rgba(59,130,246,0.12);
  --cyan:#06b6d4;
  --cyan-soft:rgba(6,182,212,0.12);
  --violet:#8b5cf6;
  --violet-soft:rgba(139,92,246,0.12);
  --gradient-brand:linear-gradient(135deg,#6366f1,#8b5cf6,#a78bfa);
  --gradient-success:linear-gradient(135deg,#10b981,#06b6d4);
  --gradient-warm:linear-gradient(135deg,#f59e0b,#ef4444);
  --gradient-cool:linear-gradient(135deg,#3b82f6,#8b5cf6);
  --shadow-sm:0 1px 3px rgba(0,0,0,0.4);
  --shadow-md:0 4px 16px rgba(0,0,0,0.5);
  --shadow-lg:0 8px 32px rgba(0,0,0,0.6);
  --shadow-glow:0 0 40px rgba(99,102,241,0.15);
  --ambient-1:rgba(99,102,241,0.07);
  --ambient-2:rgba(139,92,246,0.05);
  --auth-card-bg:rgba(17,17,25,0.8);
  --radius:10px;
  --radius-lg:14px;
  --radius-xl:18px;
  --radius-2xl:24px;
  --font:'Inter',system-ui,-apple-system,sans-serif;
  --transition:all 0.25s cubic-bezier(0.4,0,0.2,1);
  --sidebar-width:260px;
  --is-dark:1;
}

/* ===== Theme: Cyberpunk ===== */
[data-theme="cyberpunk"]{
  --bg-primary:#09090d;
  --bg-secondary:#0d0d15;
  --bg-card:#12111c;
  --bg-card-hover:#1a1928;
  --bg-input:#0f0e18;
  --bg-elevated:#1e1d30;
  --border:rgba(0,255,200,0.07);
  --border-hover:rgba(0,255,200,0.12);
  --border-active:rgba(0,255,200,0.5);
  --text-primary:#e0ffe8;
  --text-secondary:#78b8a0;
  --text-muted:#3d6858;
  --accent:#00ffc8;
  --accent-light:#5cffda;
  --accent-glow:rgba(0,255,200,0.2);
  --accent-soft:rgba(0,255,200,0.06);
  --accent-subtle:rgba(0,255,200,0.06);
  --success:#00ffc8;
  --success-soft:rgba(0,255,200,0.12);
  --warning:#ffd000;
  --warning-soft:rgba(255,208,0,0.12);
  --danger:#ff3860;
  --danger-soft:rgba(255,56,96,0.12);
  --info:#00d4ff;
  --info-soft:rgba(0,212,255,0.12);
  --cyan:#00d4ff;
  --cyan-soft:rgba(0,212,255,0.12);
  --violet:#8b5cf6;
  --violet-soft:rgba(139,92,246,0.12);
  --gradient-brand:linear-gradient(135deg,#00ffc8,#00d4ff,#5cffda);
  --gradient-success:linear-gradient(135deg,#00ffc8,#7cff6b);
  --gradient-warm:linear-gradient(135deg,#ffd000,#ff8c00);
  --gradient-cool:linear-gradient(135deg,#00ffc8,#00d4ff);
  --shadow-sm:0 1px 3px rgba(0,0,0,0.5);
  --shadow-md:0 4px 16px rgba(0,0,0,0.6);
  --shadow-lg:0 8px 32px rgba(0,0,0,0.7);
  --shadow-glow:0 0 40px rgba(0,255,200,0.1);
  --ambient-1:rgba(0,255,200,0.07);
  --ambient-2:rgba(0,212,255,0.05);
  --auth-card-bg:rgba(18,17,28,0.85);
  --is-dark:1;
}
[data-theme="cyberpunk"] .btn-primary{color:#09090d;font-weight:700}
[data-theme="cyberpunk"] .btn-submit{color:#09090d;font-weight:700}
[data-theme="cyberpunk"] .sidebar-brand-text h1{background:var(--gradient-brand);-webkit-background-clip:text;-webkit-text-fill-color:transparent}
[data-theme="cyberpunk"] .auth-card h2{background:var(--gradient-brand);-webkit-background-clip:text;-webkit-text-fill-color:transparent}

/* ===== Theme: Aurora ===== */
[data-theme="aurora"]{
  --bg-primary:#070710;
  --bg-secondary:#0b0b18;
  --bg-card:#10101e;
  --bg-card-hover:#171730;
  --bg-input:#0d0d1c;
  --bg-elevated:#1c1c3a;
  --border:rgba(168,85,247,0.08);
  --border-hover:rgba(168,85,247,0.14);
  --border-active:rgba(168,85,247,0.5);
  --text-primary:#f0e8ff;
  --text-secondary:#9b8cc0;
  --text-muted:#5c4d78;
  --accent:#a855f7;
  --accent-light:#c084fc;
  --accent-glow:rgba(168,85,247,0.22);
  --accent-soft:rgba(168,85,247,0.07);
  --accent-subtle:rgba(168,85,247,0.07);
  --success:#22d3ee;
  --success-soft:rgba(34,211,238,0.12);
  --warning:#fbbf24;
  --warning-soft:rgba(251,191,36,0.12);
  --danger:#f43f5e;
  --danger-soft:rgba(244,63,94,0.12);
  --info:#60a5fa;
  --info-soft:rgba(96,165,250,0.12);
  --cyan:#22d3ee;
  --cyan-soft:rgba(34,211,238,0.12);
  --violet:#a855f7;
  --violet-soft:rgba(168,85,247,0.12);
  --gradient-brand:linear-gradient(135deg,#a855f7,#ec4899,#c084fc);
  --gradient-success:linear-gradient(135deg,#22d3ee,#a855f7);
  --gradient-warm:linear-gradient(135deg,#fbbf24,#f43f5e);
  --gradient-cool:linear-gradient(135deg,#a855f7,#ec4899);
  --shadow-sm:0 1px 3px rgba(0,0,0,0.4);
  --shadow-md:0 4px 16px rgba(0,0,0,0.5);
  --shadow-lg:0 8px 32px rgba(0,0,0,0.6);
  --shadow-glow:0 0 40px rgba(168,85,247,0.12);
  --ambient-1:rgba(168,85,247,0.07);
  --ambient-2:rgba(236,72,153,0.05);
  --auth-card-bg:rgba(16,16,30,0.85);
  --is-dark:1;
}

/* ===== Theme: Daylight ===== */
[data-theme="daylight"]{
  --bg-primary:#f8f9fc;
  --bg-secondary:#ffffff;
  --bg-card:#ffffff;
  --bg-card-hover:#f1f3f9;
  --bg-input:#f3f4f8;
  --bg-elevated:#e2e8f0;
  --border:rgba(0,0,0,0.08);
  --border-hover:rgba(0,0,0,0.12);
  --border-active:rgba(99,102,241,0.5);
  --text-primary:#1a1a2e;
  --text-secondary:#64648a;
  --text-muted:#9e9eb8;
  --accent:#6366f1;
  --accent-light:#4f46e5;
  --accent-glow:rgba(99,102,241,0.18);
  --accent-soft:rgba(99,102,241,0.06);
  --accent-subtle:rgba(99,102,241,0.06);
  --success:#059669;
  --success-soft:rgba(5,150,105,0.1);
  --warning:#d97706;
  --warning-soft:rgba(217,119,6,0.1);
  --danger:#dc2626;
  --danger-soft:rgba(220,38,38,0.1);
  --info:#2563eb;
  --info-soft:rgba(37,99,235,0.1);
  --cyan:#0891b2;
  --cyan-soft:rgba(8,145,178,0.08);
  --violet:#7c3aed;
  --violet-soft:rgba(124,58,237,0.08);
  --gradient-brand:linear-gradient(135deg,#6366f1,#8b5cf6,#a78bfa);
  --gradient-success:linear-gradient(135deg,#059669,#0891b2);
  --gradient-warm:linear-gradient(135deg,#d97706,#dc2626);
  --gradient-cool:linear-gradient(135deg,#2563eb,#7c3aed);
  --shadow-sm:0 1px 3px rgba(0,0,0,0.06);
  --shadow-md:0 4px 16px rgba(0,0,0,0.08);
  --shadow-lg:0 8px 32px rgba(0,0,0,0.12);
  --shadow-glow:0 0 40px rgba(99,102,241,0.08);
  --ambient-1:rgba(99,102,241,0.04);
  --ambient-2:rgba(139,92,246,0.03);
  --auth-card-bg:rgba(255,255,255,0.9);
  --is-dark:0;
}
[data-theme="daylight"] .auth-logo svg{stroke:#fff}
[data-theme="daylight"] .sidebar-brand-icon svg{stroke:#fff}
[data-theme="daylight"] .user-avatar{color:#fff}
[data-theme="daylight"] .step-num{color:#fff}
[data-theme="daylight"] .btn-primary{color:#fff}
[data-theme="daylight"] .btn-submit{color:#fff}
[data-theme="daylight"] ::-webkit-scrollbar-thumb{background:rgba(0,0,0,0.12)}
[data-theme="daylight"] ::-webkit-scrollbar-thumb:hover{background:rgba(0,0,0,0.2)}

/* ===== Theme: Sunset ===== */
[data-theme="sunset"]{
  --bg-primary:#0f0a08;
  --bg-secondary:#161010;
  --bg-card:#1c1412;
  --bg-card-hover:#261c18;
  --bg-input:#140f0d;
  --bg-elevated:#2e2018;
  --border:rgba(251,146,60,0.08);
  --border-hover:rgba(251,146,60,0.14);
  --border-active:rgba(249,115,22,0.5);
  --text-primary:#fde8d8;
  --text-secondary:#b89080;
  --text-muted:#785848;
  --accent:#f97316;
  --accent-light:#fb923c;
  --accent-glow:rgba(249,115,22,0.22);
  --accent-soft:rgba(249,115,22,0.07);
  --accent-subtle:rgba(249,115,22,0.07);
  --success:#34d399;
  --success-soft:rgba(52,211,153,0.12);
  --warning:#fbbf24;
  --warning-soft:rgba(251,191,36,0.12);
  --danger:#ef4444;
  --danger-soft:rgba(239,68,68,0.12);
  --info:#60a5fa;
  --info-soft:rgba(96,165,250,0.12);
  --cyan:#06b6d4;
  --cyan-soft:rgba(6,182,212,0.12);
  --violet:#8b5cf6;
  --violet-soft:rgba(139,92,246,0.12);
  --gradient-brand:linear-gradient(135deg,#f97316,#ef4444,#fb923c);
  --gradient-success:linear-gradient(135deg,#34d399,#fbbf24);
  --gradient-warm:linear-gradient(135deg,#fbbf24,#f97316);
  --gradient-cool:linear-gradient(135deg,#f97316,#ef4444);
  --shadow-sm:0 1px 3px rgba(0,0,0,0.5);
  --shadow-md:0 4px 16px rgba(0,0,0,0.6);
  --shadow-lg:0 8px 32px rgba(0,0,0,0.7);
  --shadow-glow:0 0 40px rgba(249,115,22,0.1);
  --ambient-1:rgba(249,115,22,0.07);
  --ambient-2:rgba(239,68,68,0.05);
  --auth-card-bg:rgba(28,20,18,0.85);
  --is-dark:1;
}

body{font-family:var(--font);background:var(--bg-primary);color:var(--text-primary);line-height:1.6;overflow-x:hidden;min-height:100vh}

::-webkit-scrollbar{width:5px}
::-webkit-scrollbar-track{background:transparent}
::-webkit-scrollbar-thumb{background:rgba(255,255,255,0.08);border-radius:10px}
::-webkit-scrollbar-thumb:hover{background:rgba(255,255,255,0.15)}

/* Ambient background */
.ambient{position:fixed;inset:0;z-index:0;pointer-events:none;overflow:hidden}
.ambient::before{content:'';position:absolute;width:800px;height:800px;top:-200px;left:-100px;background:radial-gradient(circle,var(--ambient-1) 0%,transparent 70%);animation:float1 20s ease-in-out infinite}
.ambient::after{content:'';position:absolute;width:600px;height:600px;bottom:-100px;right:-100px;background:radial-gradient(circle,var(--ambient-2) 0%,transparent 70%);animation:float2 25s ease-in-out infinite}
@keyframes float1{0%,100%{transform:translate(0,0)}50%{transform:translate(60px,40px)}}
@keyframes float2{0%,100%{transform:translate(0,0)}50%{transform:translate(-40px,-60px)}}

/* SVG icon base */
.icon{width:20px;height:20px;fill:none;stroke:currentColor;stroke-width:1.8;stroke-linecap:round;stroke-linejoin:round;flex-shrink:0}
.icon-sm{width:16px;height:16px}
.icon-lg{width:24px;height:24px}

/* ===== Auth ===== */
.auth-container{display:flex;align-items:center;justify-content:center;min-height:100vh;position:relative;z-index:1;padding:20px}
.auth-card{background:var(--auth-card-bg);backdrop-filter:blur(24px);border:1px solid var(--border);border-radius:var(--radius-2xl);padding:48px 44px;width:100%;max-width:420px;box-shadow:var(--shadow-lg),var(--shadow-glow)}
.auth-brand{display:flex;flex-direction:column;align-items:center;margin-bottom:36px}
.auth-logo{width:64px;height:64px;background:var(--gradient-brand);border-radius:18px;display:flex;align-items:center;justify-content:center;box-shadow:var(--shadow-glow);margin-bottom:20px}
.auth-logo svg{width:32px;height:32px;stroke:#fff;fill:none;stroke-width:1.8;stroke-linecap:round;stroke-linejoin:round}
.auth-card h2{font-size:24px;font-weight:800;text-align:center;margin-bottom:4px;background:var(--gradient-brand);-webkit-background-clip:text;-webkit-text-fill-color:transparent}
.auth-card .subtitle{text-align:center;color:var(--text-secondary);font-size:13px}
.auth-card .form-group{margin-bottom:18px}
.auth-card .form-group label{display:block;font-size:12px;font-weight:600;color:var(--text-secondary);margin-bottom:7px;letter-spacing:0.3px}
.auth-card .form-input{width:100%;padding:11px 16px;background:var(--bg-input);border:1px solid var(--border);border-radius:var(--radius);color:var(--text-primary);font-size:14px;font-family:var(--font);transition:var(--transition)}
.auth-card .form-input:focus{outline:none;border-color:var(--accent);box-shadow:0 0 0 3px var(--accent-glow);background:var(--bg-card)}
.auth-card .form-input::placeholder{color:var(--text-muted)}
.auth-card .btn-submit{width:100%;padding:12px;font-size:14px;font-weight:700;border:none;border-radius:var(--radius);background:var(--gradient-brand);color:#fff;cursor:pointer;transition:var(--transition);margin-top:6px;font-family:var(--font);letter-spacing:0.3px;position:relative;overflow:hidden}
.auth-card .btn-submit:hover{transform:translateY(-1px);box-shadow:0 6px 24px var(--accent-glow)}
.auth-card .btn-submit:active{transform:translateY(0)}
.auth-card .btn-submit:disabled{opacity:0.5;cursor:not-allowed;transform:none}
.auth-card .switch-link{text-align:center;margin-top:24px;font-size:13px;color:var(--text-secondary)}
.auth-card .switch-link a{color:var(--accent-light);cursor:pointer;text-decoration:none;font-weight:600;transition:color .2s}
.auth-card .switch-link a:hover{color:#a78bfa}
.auth-card .error-msg{color:var(--danger);font-size:13px;margin-top:8px;text-align:center;min-height:20px}

/* ===== App Layout ===== */
.app-container{display:none;min-height:100vh;position:relative;z-index:1}

/* Sidebar */
.sidebar{width:var(--sidebar-width);background:var(--bg-secondary);border-right:1px solid var(--border);position:fixed;top:0;left:0;bottom:0;z-index:100;display:flex;flex-direction:column;transition:var(--transition)}
.sidebar-header{padding:24px 20px;border-bottom:1px solid var(--border)}
.sidebar-brand{display:flex;align-items:center;gap:12px}
.sidebar-brand-icon{width:38px;height:38px;background:var(--gradient-brand);border-radius:11px;display:flex;align-items:center;justify-content:center;box-shadow:var(--shadow-glow)}
.sidebar-brand-icon svg{width:20px;height:20px;stroke:#fff;fill:none;stroke-width:2;stroke-linecap:round;stroke-linejoin:round}
.sidebar-brand-text h1{font-size:15px;font-weight:700;background:var(--gradient-brand);-webkit-background-clip:text;-webkit-text-fill-color:transparent;line-height:1.3}
.sidebar-brand-text p{font-size:10px;color:var(--text-muted);letter-spacing:0.3px}

.sidebar-nav{flex:1;padding:16px 10px;overflow-y:auto}
.nav-section{margin-bottom:24px}
.nav-section-label{font-size:10px;font-weight:700;text-transform:uppercase;letter-spacing:1.5px;color:var(--text-muted);padding:0 14px;margin-bottom:8px}
.nav-item{display:flex;align-items:center;gap:11px;padding:10px 14px;border-radius:var(--radius);color:var(--text-secondary);cursor:pointer;transition:var(--transition);font-size:13px;font-weight:500;position:relative;border:none;background:none;width:100%;text-align:left;font-family:var(--font)}
.nav-item:hover{background:var(--accent-soft);color:var(--text-primary)}
.nav-item.active{background:var(--accent-soft);color:var(--accent-light);font-weight:600}
.nav-item.active::before{content:'';position:absolute;left:0;top:50%;transform:translateY(-50%);width:3px;height:55%;background:var(--gradient-brand);border-radius:0 4px 4px 0}
.nav-item svg{width:18px;height:18px;stroke:currentColor;fill:none;stroke-width:1.8;stroke-linecap:round;stroke-linejoin:round;flex-shrink:0}

.sidebar-footer{padding:14px 16px;border-top:1px solid var(--border)}
.sidebar-user{display:flex;align-items:center;gap:10px;padding:8px;border-radius:var(--radius);transition:var(--transition)}
.sidebar-user:hover{background:rgba(255,255,255,0.03)}
.user-avatar{width:34px;height:34px;border-radius:10px;background:var(--gradient-brand);display:flex;align-items:center;justify-content:center;font-size:14px;font-weight:700;color:#fff;flex-shrink:0}
.user-meta{flex:1;min-width:0}
.user-meta .name{font-size:13px;font-weight:600;white-space:nowrap;overflow:hidden;text-overflow:ellipsis}
.user-meta .role{font-size:11px;color:var(--text-muted)}
.btn-logout{padding:6px;background:none;border:1px solid var(--border);border-radius:8px;color:var(--text-muted);cursor:pointer;transition:var(--transition);display:flex;align-items:center;justify-content:center}
.btn-logout:hover{border-color:var(--danger);color:var(--danger);background:var(--danger-soft)}
.btn-logout svg{width:16px;height:16px;stroke:currentColor;fill:none;stroke-width:1.8;stroke-linecap:round;stroke-linejoin:round}

/* Main */
.main{margin-left:var(--sidebar-width);min-height:100vh;padding:32px 36px}

/* Page transitions */
.page{display:none;animation:pageIn 0.35s ease}
.page.active{display:block}
@keyframes pageIn{from{opacity:0;transform:translateY(10px)}to{opacity:1;transform:translateY(0)}}

/* Page header */
.page-header{margin-bottom:28px;display:flex;justify-content:space-between;align-items:flex-start}
.page-header-text h2{font-size:22px;font-weight:700;margin-bottom:4px;letter-spacing:-0.3px}
.page-header-text p{color:var(--text-secondary);font-size:13px}

/* Stat cards */
.stats-grid{display:grid;grid-template-columns:repeat(auto-fit,minmax(210px,1fr));gap:16px;margin-bottom:30px}
.stat-card{background:var(--bg-card);border:1px solid var(--border);border-radius:var(--radius-lg);padding:22px;position:relative;overflow:hidden;transition:var(--transition)}
.stat-card:hover{border-color:var(--border-hover);transform:translateY(-2px);box-shadow:var(--shadow-md)}
.stat-card-header{display:flex;justify-content:space-between;align-items:center;margin-bottom:14px}
.stat-card-icon{width:40px;height:40px;border-radius:11px;display:flex;align-items:center;justify-content:center}
.stat-card-icon svg{width:20px;height:20px;stroke:currentColor;fill:none;stroke-width:1.8;stroke-linecap:round;stroke-linejoin:round}
.stat-card-icon.purple{background:var(--violet-soft);color:var(--violet)}
.stat-card-icon.blue{background:var(--info-soft);color:var(--info)}
.stat-card-icon.cyan{background:var(--cyan-soft);color:var(--cyan)}
.stat-card-icon.green{background:var(--success-soft);color:var(--success)}
.stat-label{font-size:12px;color:var(--text-secondary);font-weight:500;letter-spacing:0.2px}
.stat-value{font-size:28px;font-weight:800;letter-spacing:-0.5px;line-height:1.2}
.stat-card::after{content:'';position:absolute;bottom:0;left:0;right:0;height:2px;opacity:0;transition:opacity .3s}
.stat-card:hover::after{opacity:1}
.stat-card:nth-child(1)::after{background:var(--gradient-brand)}
.stat-card:nth-child(2)::after{background:var(--gradient-cool)}
.stat-card:nth-child(3)::after{background:var(--gradient-success)}
.stat-card:nth-child(4)::after{background:var(--gradient-warm)}

/* Card */
.card{background:var(--bg-card);border:1px solid var(--border);border-radius:var(--radius-lg);transition:var(--transition);overflow:hidden}
.card:hover{border-color:var(--border-hover)}
.card-header{display:flex;justify-content:space-between;align-items:center;padding:18px 22px;border-bottom:1px solid var(--border)}
.card-title{font-size:15px;font-weight:600;display:flex;align-items:center;gap:8px}
.card-body{padding:0}

/* Table */
.data-table{width:100%;border-collapse:collapse}
.data-table th{padding:12px 20px;text-align:left;font-size:11px;text-transform:uppercase;letter-spacing:1px;color:var(--text-muted);font-weight:600;background:rgba(255,255,255,0.02)}
.data-table td{padding:13px 20px;font-size:13px;border-top:1px solid var(--border);color:var(--text-secondary)}
.data-table tr:hover td{background:rgba(255,255,255,0.02);color:var(--text-primary)}

/* Badges */
.badge{display:inline-flex;align-items:center;gap:5px;padding:4px 10px;border-radius:20px;font-size:11px;font-weight:600;letter-spacing:0.2px}
.badge::before{content:'';width:6px;height:6px;border-radius:50%;flex-shrink:0}
.badge-success{background:var(--success-soft);color:var(--success)}.badge-success::before{background:var(--success)}
.badge-danger{background:var(--danger-soft);color:var(--danger)}.badge-danger::before{background:var(--danger)}
.badge-info{background:var(--info-soft);color:var(--info)}.badge-info::before{background:var(--info)}
.badge-warning{background:var(--warning-soft);color:var(--warning)}.badge-warning::before{background:var(--warning)}

/* Buttons */
.btn{display:inline-flex;align-items:center;gap:7px;padding:9px 18px;border-radius:var(--radius);font-size:13px;font-weight:600;cursor:pointer;transition:var(--transition);border:none;font-family:var(--font)}
.btn-sm{padding:6px 14px;font-size:12px;border-radius:8px}
.btn-primary{background:var(--gradient-brand);color:#fff;box-shadow:0 2px 12px var(--accent-glow)}
.btn-primary:hover{transform:translateY(-1px);box-shadow:0 6px 24px var(--accent-glow)}
.btn-primary:active{transform:translateY(0)}
.btn-danger{background:var(--danger-soft);color:var(--danger);border:1px solid rgba(239,68,68,0.2)}
.btn-danger:hover{background:rgba(239,68,68,0.2);border-color:rgba(239,68,68,0.4)}
.btn-ghost{background:transparent;color:var(--text-secondary);border:1px solid var(--border)}
.btn-ghost:hover{border-color:var(--accent);color:var(--accent);background:var(--accent-soft)}

/* Modal */
.modal-overlay{position:fixed;inset:0;background:rgba(0,0,0,0.65);backdrop-filter:blur(8px);z-index:200;display:none;align-items:center;justify-content:center}
.modal-overlay.show{display:flex}
.modal{background:var(--bg-card);border:1px solid var(--border);border-radius:var(--radius-xl);padding:32px;width:90%;max-width:480px;box-shadow:var(--shadow-lg);animation:modalIn .25s ease}
@keyframes modalIn{from{opacity:0;transform:scale(0.95)}to{opacity:1;transform:scale(1)}}
.modal h3{font-size:18px;font-weight:700;margin-bottom:20px;display:flex;align-items:center;gap:10px}
.modal .form-group{margin-bottom:16px}
.modal .form-group label{display:block;font-size:12px;font-weight:600;color:var(--text-secondary);margin-bottom:6px;letter-spacing:.3px}
.modal .form-input{width:100%;padding:10px 14px;background:var(--bg-input);border:1px solid var(--border);border-radius:var(--radius);color:var(--text-primary);font-size:14px;font-family:var(--font);transition:var(--transition)}
.modal .form-input:focus{outline:none;border-color:var(--accent);box-shadow:0 0 0 3px var(--accent-glow)}
.modal .form-input::placeholder{color:var(--text-muted)}
.modal .modal-actions{display:flex;gap:10px;justify-content:flex-end;margin-top:24px}

/* Toast */
.toast{position:fixed;top:24px;right:24px;padding:14px 22px;border-radius:var(--radius);font-size:13px;font-weight:600;z-index:300;transform:translateX(130%);transition:transform .35s cubic-bezier(0.4,0,0.2,1);max-width:380px;border:1px solid transparent;backdrop-filter:blur(12px);display:flex;align-items:center;gap:10px}
.toast.show{transform:translateX(0)}
.toast-success{background:rgba(16,185,129,0.15);color:var(--success);border-color:rgba(16,185,129,0.25)}
.toast-error{background:rgba(239,68,68,0.15);color:var(--danger);border-color:rgba(239,68,68,0.25)}

/* Key display */
.key-display{background:var(--bg-primary);border:1px solid var(--border-active);border-radius:var(--radius);padding:14px 50px 14px 16px;font-family:'Cascadia Code','Fira Code',monospace;font-size:12.5px;word-break:break-all;margin:12px 0;position:relative;color:var(--accent-light);line-height:1.5}
.key-display .copy-chip{position:absolute;right:8px;top:50%;transform:translateY(-50%);padding:4px 10px;background:var(--accent);color:#fff;border:none;border-radius:6px;cursor:pointer;font-size:11px;font-weight:700;font-family:var(--font);transition:var(--transition)}
.key-display .copy-chip:hover{background:var(--accent-light)}

/* Code block */
.code-block{background:var(--bg-primary);border:1px solid var(--border);border-radius:var(--radius);padding:18px 20px;font-family:'Cascadia Code','Fira Code','Consolas',monospace;font-size:12px;line-height:1.7;overflow-x:auto;white-space:pre;color:var(--text-secondary);position:relative;margin-top:12px}
.code-block .copy-chip{position:absolute;right:10px;top:10px;padding:4px 12px;background:var(--bg-elevated);color:var(--text-secondary);border:1px solid var(--border);border-radius:6px;cursor:pointer;font-size:11px;font-family:var(--font);transition:var(--transition)}
.code-block .copy-chip:hover{border-color:var(--accent);color:var(--accent)}

/* Model grid */
.model-grid{display:grid;grid-template-columns:repeat(auto-fill,minmax(280px,1fr));gap:14px}
.model-card{background:var(--bg-card);border:1px solid var(--border);border-radius:var(--radius-lg);padding:22px;transition:var(--transition);position:relative;overflow:hidden}
.model-card:hover{border-color:var(--border-hover);transform:translateY(-2px);box-shadow:var(--shadow-md)}
.model-card::before{content:'';position:absolute;top:0;left:0;right:0;height:2px;background:var(--gradient-cool);opacity:0;transition:opacity .3s}
.model-card:hover::before{opacity:1}
.model-card-header{display:flex;justify-content:space-between;align-items:flex-start;margin-bottom:14px}
.model-card-name{font-size:15px;font-weight:700;color:var(--text-primary)}
.model-card-meta{display:flex;gap:14px;font-size:12px;color:var(--text-secondary)}
.model-card-meta span{display:flex;align-items:center;gap:5px}

/* Empty state */
.empty-state{text-align:center;padding:60px 20px;color:var(--text-muted)}
.empty-state-icon{width:56px;height:56px;margin:0 auto 16px;background:rgba(255,255,255,0.03);border-radius:16px;display:flex;align-items:center;justify-content:center}
.empty-state-icon svg{width:28px;height:28px;stroke:var(--text-muted);fill:none;stroke-width:1.5;stroke-linecap:round;stroke-linejoin:round}
.empty-state p{font-size:14px;line-height:1.7}

/* Steps */
.step-card{background:var(--bg-card);border:1px solid var(--border);border-radius:var(--radius-lg);padding:26px;margin-bottom:16px;transition:var(--transition);display:flex;gap:18px;align-items:flex-start}
.step-card:hover{border-color:var(--border-hover)}
.step-num{width:36px;height:36px;border-radius:10px;display:flex;align-items:center;justify-content:center;font-size:14px;font-weight:800;flex-shrink:0;color:#fff}
.step-num.s1{background:var(--gradient-brand)}
.step-num.s2{background:var(--gradient-cool)}
.step-num.s3{background:var(--gradient-success)}
.step-content h3{font-size:15px;font-weight:700;margin-bottom:6px}
.step-content p{color:var(--text-secondary);font-size:13px;line-height:1.6}
.step-content a{color:var(--accent-light);cursor:pointer;text-decoration:none;font-weight:600}
.step-content a:hover{text-decoration:underline}
.step-content code{background:var(--bg-primary);padding:2px 8px;border-radius:5px;color:var(--accent-light);font-size:12px}

/* Theme Switcher FAB */
.theme-fab{position:fixed;bottom:28px;right:28px;z-index:150;width:46px;height:46px;border-radius:50%;background:var(--gradient-brand);border:none;cursor:pointer;display:flex;align-items:center;justify-content:center;font-size:20px;box-shadow:0 4px 20px var(--accent-glow),0 0 0 1px rgba(255,255,255,0.06);transition:all 0.3s cubic-bezier(0.34,1.56,0.64,1);color:#fff}
.theme-fab:hover{transform:scale(1.1);box-shadow:0 6px 28px var(--accent-glow)}
.theme-fab.open{transform:rotate(45deg) scale(1.05)}
.theme-panel{position:fixed;bottom:86px;right:28px;z-index:150;background:var(--bg-card);border:1px solid var(--border-hover);border-radius:var(--radius-lg);padding:8px;width:230px;box-shadow:var(--shadow-lg);opacity:0;transform:translateY(12px) scale(0.95);pointer-events:none;transition:all 0.25s cubic-bezier(0.34,1.56,0.64,1)}
.theme-panel.open{opacity:1;transform:translateY(0) scale(1);pointer-events:auto}
.theme-panel-title{font-size:10px;font-weight:700;text-transform:uppercase;letter-spacing:1.8px;color:var(--text-muted);padding:10px 12px 6px}
.theme-option{display:flex;align-items:center;gap:12px;padding:10px 12px;border-radius:var(--radius);cursor:pointer;transition:all 0.15s ease;border:1px solid transparent;margin-bottom:2px}
.theme-option:last-child{margin-bottom:0}
.theme-option:hover{background:var(--accent-soft);border-color:var(--border)}
.theme-option.active{background:var(--accent-soft);border-color:var(--border-active)}
.theme-option-colors{display:flex;gap:3px;flex-shrink:0}
.theme-option-dot{width:14px;height:14px;border-radius:50%;box-shadow:inset 0 1px 2px rgba(0,0,0,0.2)}
.theme-option-info{flex:1;min-width:0}
.theme-option-name{font-size:13px;font-weight:600;color:var(--text-primary);line-height:1.2}
.theme-option-desc{font-size:10px;color:var(--text-muted);margin-top:2px}

/* Top Bar (mobile) */
.top-bar{display:none;align-items:center;justify-content:space-between;padding:0 0 20px;gap:16px}
.top-bar-left{display:flex;align-items:center;gap:12px}
.mobile-menu-btn{display:none;align-items:center;justify-content:center;width:42px;height:42px;border-radius:var(--radius);background:var(--bg-card);border:1px solid var(--border);color:var(--text-secondary);cursor:pointer;transition:var(--transition);flex-shrink:0}
.mobile-menu-btn:hover{background:var(--bg-card-hover);color:var(--text-primary);border-color:var(--border-hover)}
.mobile-menu-btn svg{width:20px;height:20px;stroke:currentColor;fill:none;stroke-width:2;stroke-linecap:round;stroke-linejoin:round}
.sidebar-overlay{display:none;position:fixed;inset:0;background:rgba(0,0,0,0.6);backdrop-filter:blur(4px);z-index:99;opacity:0;transition:opacity .3s}
.sidebar-overlay.show{display:block;opacity:1}
.sidebar-close{display:none;position:absolute;top:16px;right:16px;width:32px;height:32px;border-radius:8px;background:rgba(255,255,255,0.06);border:1px solid var(--border);color:var(--text-secondary);cursor:pointer;align-items:center;justify-content:center;transition:var(--transition);z-index:2}
.sidebar-close:hover{background:var(--danger-soft);color:var(--danger);border-color:rgba(239,68,68,0.3)}
.sidebar-close svg{width:16px;height:16px;stroke:currentColor;fill:none;stroke-width:2;stroke-linecap:round;stroke-linejoin:round}

/* Action Buttons */
.action-group{display:flex;gap:6px;align-items:center}
.btn-icon{display:inline-flex;align-items:center;justify-content:center;width:32px;height:32px;border-radius:8px;border:1px solid var(--border);background:transparent;color:var(--text-secondary);cursor:pointer;transition:var(--transition)}
.btn-icon:hover{border-color:var(--accent);color:var(--accent);background:var(--accent-soft)}
.btn-icon.btn-icon-danger:hover{border-color:var(--danger);color:var(--danger);background:var(--danger-soft)}
.btn-icon svg{width:15px;height:15px;stroke:currentColor;fill:none;stroke-width:2;stroke-linecap:round;stroke-linejoin:round}

/* Responsive */
@media(max-width:1024px){
  .main{padding:24px 20px}
  .stats-grid{grid-template-columns:repeat(2,1fr)}
}
@media(max-width:768px){
  .sidebar{transform:translateX(-100%);z-index:101;box-shadow:var(--shadow-lg)}
  .sidebar.open{transform:translateX(0)}
  .sidebar-close{display:flex}
  .main{margin-left:0;padding:20px 16px}
  .top-bar{display:flex}
  .mobile-menu-btn{display:flex}
  .stats-grid{grid-template-columns:1fr 1fr}
  .auth-card{padding:36px 24px}
  .model-grid{grid-template-columns:1fr}
  .page-header{flex-direction:column;gap:16px;align-items:stretch}
  .page-header .btn{width:100%;justify-content:center}
  .data-table{font-size:12px}
  .data-table th,.data-table td{padding:10px 12px}
  .modal{width:95%;padding:24px;max-width:none}
  .card-header{padding:14px 16px}
  .step-card{padding:20px;gap:14px}
}
@media(max-width:480px){
  .stats-grid{grid-template-columns:1fr}
  .stat-value{font-size:24px}
}
</style>
</head>
<body>
<div class="ambient"></div>

<!-- Auth Views -->
<div id="authView" class="auth-container">
  <div class="auth-card" id="loginCard">
    <div class="auth-brand">
      <div class="auth-logo">
        <svg viewBox="0 0 24 24"><path d="M12 2L2 7l10 5 10-5-10-5z"/><path d="M2 17l10 5 10-5"/><path d="M2 12l10 5 10-5"/></svg>
      </div>
      <h2>AI Gateway</h2>
      <p class="subtitle">统一 AI 模型 API 网关 · 用户中心</p>
    </div>
    <div class="form-group">
      <label>用户名</label>
      <input type="text" class="form-input" id="loginUser" placeholder="请输入用户名" autocomplete="username">
    </div>
    <div class="form-group">
      <label>密码</label>
      <input type="password" class="form-input" id="loginPass" placeholder="请输入密码" autocomplete="current-password" onkeypress="if(event.key==='Enter')doLogin()">
    </div>
    <div class="error-msg" id="loginError"></div>
    <button class="btn-submit" onclick="doLogin()" id="loginBtn">登 录</button>
    <div class="switch-link">还没有账号？<a onclick="showRegister()">立即注册</a></div>
  </div>

  <div class="auth-card" id="registerCard" style="display:none">
    <div class="auth-brand">
      <div class="auth-logo">
        <svg viewBox="0 0 24 24"><path d="M16 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/><circle cx="8.5" cy="7" r="4"/><line x1="20" y1="8" x2="20" y2="14"/><line x1="23" y1="11" x2="17" y2="11"/></svg>
      </div>
      <h2>创建账号</h2>
      <p class="subtitle">注册后即可获取专属 API Key</p>
    </div>
    <div class="form-group">
      <label>用户名</label>
      <input type="text" class="form-input" id="regUser" placeholder="3-50个字符" autocomplete="username">
    </div>
    <div class="form-group">
      <label>邮箱</label>
      <input type="email" class="form-input" id="regEmail" placeholder="your@email.com" autocomplete="email">
    </div>
    <div class="form-group">
      <label>密码</label>
      <input type="password" class="form-input" id="regPass" placeholder="至少6位" autocomplete="new-password">
    </div>
    <div class="form-group">
      <label>确认密码</label>
      <input type="password" class="form-input" id="regPassConfirm" placeholder="再次输入密码" autocomplete="new-password" onkeypress="if(event.key==='Enter')doRegister()">
    </div>
    <div class="error-msg" id="regError"></div>
    <button class="btn-submit" onclick="doRegister()" id="regBtn">注 册</button>
    <div class="switch-link">已有账号？<a onclick="showLogin()">返回登录</a></div>
  </div>
</div>

<!-- Main App -->
<div class="app-container" id="appView">
  <!-- Sidebar -->
  <aside class="sidebar" id="sidebar">
    <button class="sidebar-close" onclick="closeSidebar()">
      <svg viewBox="0 0 24 24"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
    </button>
    <div class="sidebar-header">
      <div class="sidebar-brand">
        <div class="sidebar-brand-icon">
          <svg viewBox="0 0 24 24"><path d="M12 2L2 7l10 5 10-5-10-5z"/><path d="M2 17l10 5 10-5"/><path d="M2 12l10 5 10-5"/></svg>
        </div>
        <div class="sidebar-brand-text">
          <h1>AI Gateway</h1>
          <p>用户中心</p>
        </div>
      </div>
    </div>
    <nav class="sidebar-nav">
      <div class="nav-section">
        <div class="nav-section-label">概览</div>
        <button class="nav-item active" data-page="dashboard" onclick="showPage('dashboard',this)">
          <svg viewBox="0 0 24 24"><rect x="3" y="3" width="7" height="7" rx="1"/><rect x="14" y="3" width="7" height="7" rx="1"/><rect x="3" y="14" width="7" height="7" rx="1"/><rect x="14" y="14" width="7" height="7" rx="1"/></svg>
          总览
        </button>
      </div>
      <div class="nav-section">
        <div class="nav-section-label">管理</div>
        <button class="nav-item" data-page="apikeys" onclick="showPage('apikeys',this)">
          <svg viewBox="0 0 24 24"><path d="M21 2l-2 2m-7.61 7.61a5.5 5.5 0 1 1-7.778 7.778 5.5 5.5 0 0 1 7.777-7.777zm0 0L15.5 7.5m0 0l3 3L22 7l-3-3m-3.5 3.5L19 4"/></svg>
          API Keys
        </button>
        <button class="nav-item" data-page="models" onclick="showPage('models',this)">
          <svg viewBox="0 0 24 24"><path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"/><polyline points="3.27 6.96 12 12.01 20.73 6.96"/><line x1="12" y1="22.08" x2="12" y2="12"/></svg>
          模型列表
        </button>
      </div>
      <div class="nav-section">
        <div class="nav-section-label">数据</div>
        <button class="nav-item" data-page="usage" onclick="showPage('usage',this)">
          <svg viewBox="0 0 24 24"><line x1="18" y1="20" x2="18" y2="10"/><line x1="12" y1="20" x2="12" y2="4"/><line x1="6" y1="20" x2="6" y2="14"/></svg>
          用量统计
        </button>
      </div>
      <div class="nav-section">
        <div class="nav-section-label">帮助</div>
        <button class="nav-item" data-page="quickstart" onclick="showPage('quickstart',this)">
          <svg viewBox="0 0 24 24"><polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2"/></svg>
          快速开始
        </button>
      </div>
    </nav>
    <div class="sidebar-footer">
      <div class="sidebar-user">
        <div class="user-avatar" id="userAvatar">U</div>
        <div class="user-meta">
          <div class="name" id="usernameDisplay">user</div>
          <div class="role">用户</div>
        </div>
        <button class="btn-logout" onclick="doLogout()" title="退出登录">
          <svg viewBox="0 0 24 24"><path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/><polyline points="16 17 21 12 16 7"/><line x1="21" y1="12" x2="9" y2="12"/></svg>
        </button>
      </div>
    </div>
  </aside>

  <div class="sidebar-overlay" id="sidebarOverlay" onclick="closeSidebar()"></div>

  <!-- Main Content -->
  <div class="main">
    <!-- Mobile Top Bar -->
    <div class="top-bar">
      <div class="top-bar-left">
        <button class="mobile-menu-btn" onclick="openSidebar()">
          <svg viewBox="0 0 24 24"><line x1="3" y1="6" x2="21" y2="6"/><line x1="3" y1="12" x2="21" y2="12"/><line x1="3" y1="18" x2="21" y2="18"/></svg>
        </button>
        <div class="sidebar-brand">
          <div class="sidebar-brand-icon" style="width:34px;height:34px;border-radius:9px">
            <svg viewBox="0 0 24 24" style="width:17px;height:17px;stroke:#fff;fill:none;stroke-width:2;stroke-linecap:round;stroke-linejoin:round"><path d="M12 2L2 7l10 5 10-5-10-5z"/><path d="M2 17l10 5 10-5"/><path d="M2 12l10 5 10-5"/></svg>
          </div>
          <div class="sidebar-brand-text">
            <h1 style="font-size:14px">AI Gateway</h1>
          </div>
        </div>
      </div>
    </div>

    <!-- Dashboard -->
    <div class="page active" id="page-dashboard">
      <div class="page-header">
        <div class="page-header-text">
          <h2>总览</h2>
          <p>欢迎回来！这里是您的 API 使用概况</p>
        </div>
      </div>
      <div class="stats-grid">
        <div class="stat-card">
          <div class="stat-card-header">
            <div class="stat-label">API Keys</div>
            <div class="stat-card-icon purple">
              <svg viewBox="0 0 24 24"><path d="M21 2l-2 2m-7.61 7.61a5.5 5.5 0 1 1-7.778 7.778 5.5 5.5 0 0 1 7.777-7.777zm0 0L15.5 7.5m0 0l3 3L22 7l-3-3m-3.5 3.5L19 4"/></svg>
            </div>
          </div>
          <div class="stat-value" id="statKeys">-</div>
        </div>
        <div class="stat-card">
          <div class="stat-card-header">
            <div class="stat-label">总调用次数</div>
            <div class="stat-card-icon blue">
              <svg viewBox="0 0 24 24"><line x1="18" y1="20" x2="18" y2="10"/><line x1="12" y1="20" x2="12" y2="4"/><line x1="6" y1="20" x2="6" y2="14"/></svg>
            </div>
          </div>
          <div class="stat-value" id="statCalls">-</div>
        </div>
        <div class="stat-card">
          <div class="stat-card-header">
            <div class="stat-label">总 Token 消耗</div>
            <div class="stat-card-icon cyan">
              <svg viewBox="0 0 24 24"><circle cx="12" cy="12" r="10"/><path d="M16 8h-6a2 2 0 1 0 0 4h4a2 2 0 1 1 0 4H8"/><path d="M12 18V6"/></svg>
            </div>
          </div>
          <div class="stat-value" id="statTokens">-</div>
        </div>
        <div class="stat-card">
          <div class="stat-card-header">
            <div class="stat-label">可用模型数</div>
            <div class="stat-card-icon green">
              <svg viewBox="0 0 24 24"><path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"/></svg>
            </div>
          </div>
          <div class="stat-value" id="statModels">-</div>
        </div>
      </div>
      <div class="card">
        <div class="card-header">
          <div class="card-title">
            <svg class="icon icon-sm" viewBox="0 0 24 24"><polyline points="22 12 18 12 15 21 9 3 6 12 2 12"/></svg>
            最近调用
          </div>
        </div>
        <div class="card-body" id="recentUsage"></div>
      </div>
    </div>

    <!-- API Keys -->
    <div class="page" id="page-apikeys">
      <div class="page-header">
        <div class="page-header-text">
          <h2>API Keys</h2>
          <p>管理您的 API 访问密钥</p>
        </div>
        <button class="btn btn-primary" onclick="showCreateKeyModal()">
          <svg class="icon icon-sm" viewBox="0 0 24 24"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
          创建 Key
        </button>
      </div>
      <div class="card">
        <div class="card-body" id="keysList"></div>
      </div>
    </div>

    <!-- Models -->
    <div class="page" id="page-models">
      <div class="page-header">
        <div class="page-header-text">
          <h2>可用模型</h2>
          <p>当前网关支持的所有 AI 模型</p>
        </div>
      </div>
      <div id="modelsList"></div>
    </div>

    <!-- Usage -->
    <div class="page" id="page-usage">
      <div class="page-header">
        <div class="page-header-text">
          <h2>用量统计</h2>
          <p>查看 API 调用详情和 Token 消耗</p>
        </div>
      </div>
      <div class="stats-grid" id="usageStats"></div>
      <div class="card">
        <div class="card-header">
          <div class="card-title">
            <svg class="icon icon-sm" viewBox="0 0 24 24"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/><line x1="16" y1="13" x2="8" y2="13"/><line x1="16" y1="17" x2="8" y2="17"/></svg>
            调用明细
          </div>
        </div>
        <div class="card-body" id="usageDetails"></div>
      </div>
    </div>

    <!-- Quick Start -->
    <div class="page" id="page-quickstart">
      <div class="page-header">
        <div class="page-header-text">
          <h2>快速开始</h2>
          <p>几分钟内开始使用 AI Gateway API</p>
        </div>
      </div>

      <div class="step-card">
        <div class="step-num s1">1</div>
        <div class="step-content">
          <h3>获取 API Key</h3>
          <p>在 <a onclick="showPage('apikeys',document.querySelector('[data-page=apikeys]'))">API Keys</a> 页面创建一个新的密钥，创建时请妥善保存</p>
        </div>
      </div>

      <div class="step-card">
        <div class="step-num s2">2</div>
        <div class="step-content">
          <h3>调用 API</h3>
          <p>完全兼容 OpenAI 格式，可直接使用 OpenAI SDK</p>
          <p style="margin-top:4px">Base URL: <code>${location.origin}/v1</code></p>
          <div style="margin-top:16px">
            <div style="font-size:12px;font-weight:600;color:var(--text-muted);text-transform:uppercase;letter-spacing:1px;margin-bottom:6px">cURL 示例</div>
            <div class="code-block" id="curlExample"><button class="copy-chip" onclick="copyCode('curlExample')">复制</button>curl ${location.origin}/v1/chat/completions \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
  "model": "gpt-4o",
  "messages": [{"role": "user", "content": "Hello!"}]
}'</div>
          </div>
          <div style="margin-top:20px">
            <div style="font-size:12px;font-weight:600;color:var(--text-muted);text-transform:uppercase;letter-spacing:1px;margin-bottom:6px">Python (OpenAI SDK)</div>
            <div class="code-block" id="pythonExample"><button class="copy-chip" onclick="copyCode('pythonExample')">复制</button>from openai import OpenAI

client = OpenAI(
    api_key="YOUR_API_KEY",
    base_url="${location.origin}/v1"
)

response = client.chat.completions.create(
    model="gpt-4o",
    messages=[{"role": "user", "content": "Hello!"}]
)
print(response.choices[0].message.content)</div>
          </div>
        </div>
      </div>

      <div class="step-card">
        <div class="step-num s3">3</div>
        <div class="step-content">
          <h3>查看模型</h3>
          <p>在 <a onclick="showPage('models',document.querySelector('[data-page=models]'))">模型列表</a> 查看所有可用模型及其详情</p>
        </div>
      </div>
    </div>
  </div>
</div>

<!-- Create Key Modal -->
<div class="modal-overlay" id="createKeyModal">
  <div class="modal">
    <h3>
      <svg class="icon" viewBox="0 0 24 24" style="color:var(--accent)"><path d="M21 2l-2 2m-7.61 7.61a5.5 5.5 0 1 1-7.778 7.778 5.5 5.5 0 0 1 7.777-7.777zm0 0L15.5 7.5m0 0l3 3L22 7l-3-3m-3.5 3.5L19 4"/></svg>
      创建 API Key
    </h3>
    <div class="form-group">
      <label>名称</label>
      <input type="text" class="form-input" id="newKeyName" placeholder="例如: my-app-key">
    </div>
    <div class="form-group">
      <label>每分钟请求限制</label>
      <input type="number" class="form-input" id="newKeyRate" value="60" min="1" max="1000">
    </div>
    <div class="form-group">
      <label>Token 配额 (0=无限)</label>
      <input type="number" class="form-input" id="newKeyQuota" value="0" min="0">
    </div>
    <div class="modal-actions">
      <button class="btn btn-ghost" onclick="closeModal('createKeyModal')">取消</button>
      <button class="btn btn-primary" onclick="createKey()">创建</button>
    </div>
  </div>
</div>

<!-- Key Created Modal -->
<div class="modal-overlay" id="keyCreatedModal">
  <div class="modal">
    <h3>
      <svg class="icon" viewBox="0 0 24 24" style="color:var(--success)"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/></svg>
      Key 创建成功
    </h3>
    <p style="color:var(--danger);font-size:13px;margin-bottom:8px;display:flex;align-items:center;gap:6px">
      <svg class="icon icon-sm" viewBox="0 0 24 24" style="color:var(--warning)"><path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"/><line x1="12" y1="9" x2="12" y2="13"/><line x1="12" y1="17" x2="12.01" y2="17"/></svg>
      请立即复制保存，此密钥仅显示一次！
    </p>
    <div class="key-display" id="newKeyDisplay">
      <button class="copy-chip" onclick="copyNewKey()">复制</button>
      <span id="newKeyValue"></span>
    </div>
    <div class="modal-actions">
      <button class="btn btn-primary" onclick="closeModal('keyCreatedModal');loadKeys()">我已保存</button>
    </div>
  </div>
</div>

<!-- Theme FAB -->
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

<!-- Toast -->
<div class="toast" id="toast"></div>

<script>
const API = '';
let userToken = localStorage.getItem('user_token');
let currentUser = JSON.parse(localStorage.getItem('user_info') || 'null');

function esc(s) {
  const d = document.createElement('div');
  d.textContent = s;
  return d.innerHTML;
}

function showToast(msg, type='success') {
  const t = document.getElementById('toast');
  t.innerHTML = (type === 'success' ? '<svg class="icon icon-sm" viewBox="0 0 24 24"><polyline points="20 6 9 17 4 12"/></svg>' : '<svg class="icon icon-sm" viewBox="0 0 24 24"><circle cx="12" cy="12" r="10"/><line x1="15" y1="9" x2="9" y2="15"/><line x1="9" y1="9" x2="15" y2="15"/></svg>') + esc(msg);
  t.className = 'toast toast-' + type + ' show';
  setTimeout(() => t.classList.remove('show'), 3000);
}

function userHeaders() {
  return { 'Content-Type': 'application/json', 'Authorization': 'Bearer ' + userToken };
}

// ===== Auth =====
function showLogin() {
  document.getElementById('loginCard').style.display = 'block';
  document.getElementById('registerCard').style.display = 'none';
  document.getElementById('loginError').textContent = '';
}

function showRegister() {
  document.getElementById('loginCard').style.display = 'none';
  document.getElementById('registerCard').style.display = 'block';
  document.getElementById('regError').textContent = '';
}

async function doLogin() {
  const user = document.getElementById('loginUser').value.trim();
  const pass = document.getElementById('loginPass').value;
  if (!user || !pass) { document.getElementById('loginError').textContent = '请输入用户名和密码'; return; }

  const btn = document.getElementById('loginBtn');
  btn.disabled = true;
  btn.textContent = '登录中...';

  try {
    const r = await fetch(API + '/user/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username: user, password: pass })
    });
    const d = await r.json();
    if (d.code === 0 && d.data) {
      userToken = d.data.token;
      currentUser = { id: d.data.user_id, username: d.data.username, role: d.data.role };
      localStorage.setItem('user_token', userToken);
      localStorage.setItem('user_info', JSON.stringify(currentUser));
      enterApp();
    } else {
      document.getElementById('loginError').textContent = (d.error && d.error.message) || '登录失败';
    }
  } catch(e) {
    document.getElementById('loginError').textContent = '网络错误，请重试';
  } finally {
    btn.disabled = false;
    btn.textContent = '登 录';
  }
}

async function doRegister() {
  const user = document.getElementById('regUser').value.trim();
  const email = document.getElementById('regEmail').value.trim();
  const pass = document.getElementById('regPass').value;
  const pass2 = document.getElementById('regPassConfirm').value;

  if (!user || !email || !pass) { document.getElementById('regError').textContent = '请填写所有字段'; return; }
  if (pass !== pass2) { document.getElementById('regError').textContent = '两次输入的密码不一致'; return; }
  if (pass.length < 6) { document.getElementById('regError').textContent = '密码至少6位'; return; }

  const btn = document.getElementById('regBtn');
  btn.disabled = true;
  btn.textContent = '注册中...';

  try {
    const r = await fetch(API + '/user/register', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username: user, email: email, password: pass })
    });
    const d = await r.json();
    if (d.code === 0 && d.data) {
      showToast('注册成功！请登录');
      showLogin();
      document.getElementById('loginUser').value = user;
      document.getElementById('loginPass').focus();
    } else {
      document.getElementById('regError').textContent = (d.error && d.error.message) || '注册失败';
    }
  } catch(e) {
    document.getElementById('regError').textContent = '网络错误，请重试';
  } finally {
    btn.disabled = false;
    btn.textContent = '注 册';
  }
}

function doLogout() {
  userToken = null;
  currentUser = null;
  localStorage.removeItem('user_token');
  localStorage.removeItem('user_info');
  document.getElementById('appView').style.display = 'none';
  document.getElementById('authView').style.display = 'flex';
  showLogin();
}

function enterApp() {
  document.getElementById('authView').style.display = 'none';
  document.getElementById('appView').style.display = 'block';
  document.getElementById('usernameDisplay').textContent = currentUser.username;
  document.getElementById('userAvatar').textContent = currentUser.username.charAt(0).toUpperCase();
  loadDashboard();
}

// ===== Navigation =====
function showPage(name, el) {
  document.querySelectorAll('.page').forEach(p => p.classList.remove('active'));
  document.querySelectorAll('.nav-item').forEach(t => t.classList.remove('active'));
  document.getElementById('page-'+name).classList.add('active');
  if (el) el.classList.add('active');

  if (name === 'dashboard') loadDashboard();
  else if (name === 'apikeys') loadKeys();
  else if (name === 'models') loadModels();
  else if (name === 'usage') loadUsage();
}

// ===== Dashboard =====
async function loadDashboard() {
  try {
    const r = await fetch(API + '/user/api/keys', { headers: userHeaders() });
    const d = await r.json();
    if (d.code === 0) {
      const keys = d.data || [];
      document.getElementById('statKeys').textContent = keys.length;
    }
  } catch(e) {}

  try {
    const r = await fetch(API + '/user/api/usage', { headers: userHeaders() });
    const d = await r.json();
    if (d.code === 0 && d.data) {
      document.getElementById('statCalls').textContent = (d.data.total_calls || 0).toLocaleString();
      document.getElementById('statTokens').textContent = (d.data.total_tokens || 0).toLocaleString();
    }
  } catch(e) {}

  try {
    const r = await fetch(API + '/user/api/models', { headers: userHeaders() });
    const d = await r.json();
    const mlist = d.data || [];
    document.getElementById('statModels').textContent = mlist.length || 0;
  } catch(e) {}

  try {
    const r = await fetch(API + '/user/api/usage/details?limit=5', { headers: userHeaders() });
    const d = await r.json();
    if (d.code === 0 && d.data && d.data.length > 0) {
      let html = '<table class="data-table"><thead><tr><th>时间</th><th>模型</th><th>Tokens</th><th>延迟</th><th>状态</th></tr></thead><tbody>';
      d.data.forEach(u => {
        const time = new Date(u.created_at).toLocaleString('zh-CN');
        const status = u.status === 'success' ? '<span class="badge badge-success">成功</span>' : '<span class="badge badge-danger">失败</span>';
        html += '<tr><td>' + esc(time) + '</td><td><strong style="color:var(--text-primary)">' + esc(u.model_name) + '</strong></td><td>' + (u.total_tokens||0).toLocaleString() + '</td><td>' + (u.latency_ms||0) + 'ms</td><td>' + status + '</td></tr>';
      });
      html += '</tbody></table>';
      document.getElementById('recentUsage').innerHTML = html;
    } else {
      document.getElementById('recentUsage').innerHTML = '<div class="empty-state"><div class="empty-state-icon"><svg viewBox="0 0 24 24"><path d="M22 12h-4l-3 9L9 3l-3 9H2"/></svg></div><p>暂无调用记录<br><span style="color:var(--text-muted);font-size:13px">创建 API Key 后开始使用吧</span></p></div>';
    }
  } catch(e) {
    document.getElementById('recentUsage').innerHTML = '<div class="empty-state"><div class="empty-state-icon"><svg viewBox="0 0 24 24"><path d="M22 12h-4l-3 9L9 3l-3 9H2"/></svg></div><p>暂无调用记录</p></div>';
  }
}

// ===== API Keys =====
async function loadKeys() {
  try {
    const r = await fetch(API + '/user/api/keys', { headers: userHeaders() });
    const d = await r.json();
    if (d.code === 0) {
      const keys = d.data || [];
      if (keys.length === 0) {
        document.getElementById('keysList').innerHTML = '<div class="empty-state"><div class="empty-state-icon"><svg viewBox="0 0 24 24"><path d="M21 2l-2 2m-7.61 7.61a5.5 5.5 0 1 1-7.778 7.778 5.5 5.5 0 0 1 7.777-7.777zm0 0L15.5 7.5m0 0l3 3L22 7l-3-3m-3.5 3.5L19 4"/></svg></div><p>您还没有 API Key<br><span style="color:var(--text-muted);font-size:13px">点击上方按钮创建第一个</span></p></div>';
        return;
      }
      let html = '<table class="data-table"><thead><tr><th>名称</th><th>Key 前缀</th><th>速率限制</th><th>配额</th><th>状态</th><th>创建时间</th><th>操作</th></tr></thead><tbody>';
      keys.forEach(k => {
        const enabled = k.enabled ? '<span class="badge badge-success">启用</span>' : '<span class="badge badge-danger">禁用</span>';
        const quota = k.quota_limit > 0 ? ((k.quota_used||0).toLocaleString() + ' / ' + k.quota_limit.toLocaleString()) : '无限';
        const time = new Date(k.created_at).toLocaleDateString('zh-CN');
        html += '<tr><td><strong style="color:var(--text-primary)">' + esc(k.name||'未命名') + '</strong></td><td><code style="background:var(--bg-primary);padding:3px 8px;border-radius:5px;font-size:12px;color:var(--accent-light)">' + esc(k.key_prefix) + '...</code></td><td>' + k.rate_limit + '/min</td><td>' + quota + '</td><td>' + enabled + '</td><td>' + esc(time) + '</td><td><button class="btn btn-danger btn-sm" onclick="deleteKey(' + k.id + ')">删除</button></td></tr>';
      });
      html += '</tbody></table>';
      document.getElementById('keysList').innerHTML = html;
    }
  } catch(e) {
    document.getElementById('keysList').innerHTML = '<div class="empty-state"><div class="empty-state-icon"><svg viewBox="0 0 24 24"><circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/></svg></div><p>加载失败，请刷新重试</p></div>';
  }
}

function showCreateKeyModal() {
  document.getElementById('newKeyName').value = '';
  document.getElementById('newKeyRate').value = '60';
  document.getElementById('newKeyQuota').value = '0';
  document.getElementById('createKeyModal').classList.add('show');
}

function closeModal(id) {
  document.getElementById(id).classList.remove('show');
}

async function createKey() {
  const name = document.getElementById('newKeyName').value.trim();
  const rate = parseInt(document.getElementById('newKeyRate').value) || 60;
  const quota = parseInt(document.getElementById('newKeyQuota').value) || 0;

  if (!name) { showToast('请输入 Key 名称', 'error'); return; }

  try {
    const r = await fetch(API + '/user/api/keys', {
      method: 'POST',
      headers: userHeaders(),
      body: JSON.stringify({ name, rate_limit: rate, quota_limit: quota })
    });
    const d = await r.json();
    if (d.code === 0 && d.data) {
      closeModal('createKeyModal');
      document.getElementById('newKeyValue').textContent = d.data.key;
      document.getElementById('keyCreatedModal').classList.add('show');
    } else {
      showToast((d.error && d.error.message) || '创建失败', 'error');
    }
  } catch(e) {
    showToast('创建失败', 'error');
  }
}

async function deleteKey(id) {
  if (!confirm('确定要删除此 API Key？此操作不可恢复。')) return;
  try {
    const r = await fetch(API + '/user/api/keys/' + id, { method: 'DELETE', headers: userHeaders() });
    const d = await r.json();
    if (d.code === 0) {
      showToast('已删除');
      loadKeys();
    } else {
      showToast((d.error && d.error.message) || '删除失败', 'error');
    }
  } catch(e) {
    showToast('删除失败', 'error');
  }
}

function copyNewKey() {
  const key = document.getElementById('newKeyValue').textContent;
  navigator.clipboard.writeText(key).then(() => showToast('已复制到剪贴板'));
}

function copyCode(id) {
  const el = document.getElementById(id);
  const text = el.textContent.replace('复制', '').trim();
  navigator.clipboard.writeText(text).then(() => showToast('已复制'));
}

// ===== Models =====
async function loadModels() {
  try {
    const r = await fetch(API + '/user/api/models', { headers: userHeaders() });
    const d = await r.json();
    const models = d.data || d.data || [];
    if (models && models.length > 0) {
      let html = '<div class="model-grid">';
      models.forEach(m => {
        const name = m.id || m.model_name || m.name || 'unknown';
        const provider = m.owned_by || m.provider_type || 'unknown';
        const mtype = m.model_type || 'chat';
        const channels = m.channel_count || 1;
        html += '<div class="model-card">' +
          '<div class="model-card-header">' +
          '<div class="model-card-name">' + esc(name) + '</div>' +
          '<span class="badge badge-info">' + esc(provider) + '</span>' +
          '</div>' +
          '<div class="model-card-meta">' +
          '<span><svg class="icon icon-sm" viewBox="0 0 24 24"><circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06A1.65 1.65 0 0 0 4.68 15a1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06A1.65 1.65 0 0 0 9 4.68a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06A1.65 1.65 0 0 0 19.4 9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"/></svg>' + esc(mtype) + '</span>' +
          '<span><svg class="icon icon-sm" viewBox="0 0 24 24"><path d="M2 20h.01"/><path d="M7 20v-4"/><path d="M12 20v-8"/><path d="M17 20V8"/><path d="M22 4v16"/></svg>' + channels + ' 通道</span>' +
          '</div></div>';
      });
      html += '</div>';
      document.getElementById('modelsList').innerHTML = html;
    } else {
      document.getElementById('modelsList').innerHTML = '<div class="empty-state"><div class="empty-state-icon"><svg viewBox="0 0 24 24"><path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"/></svg></div><p>暂无可用模型<br><span style="color:var(--text-muted);font-size:13px">请联系管理员配置</span></p></div>';
    }
  } catch(e) {
    document.getElementById('modelsList').innerHTML = '<div class="empty-state"><div class="empty-state-icon"><svg viewBox="0 0 24 24"><circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/></svg></div><p>加载失败</p></div>';
  }
}

// ===== Usage =====
async function loadUsage() {
  try {
    const r = await fetch(API + '/user/api/usage', { headers: userHeaders() });
    const d = await r.json();
    if (d.code === 0 && d.data) {
      const u = d.data;
      document.getElementById('usageStats').innerHTML =
        '<div class="stat-card"><div class="stat-card-header"><div class="stat-label">总调用</div><div class="stat-card-icon blue"><svg viewBox="0 0 24 24"><line x1="18" y1="20" x2="18" y2="10"/><line x1="12" y1="20" x2="12" y2="4"/><line x1="6" y1="20" x2="6" y2="14"/></svg></div></div><div class="stat-value">' + (u.total_calls||0).toLocaleString() + '</div></div>' +
        '<div class="stat-card"><div class="stat-card-header"><div class="stat-label">输入 Tokens</div><div class="stat-card-icon purple"><svg viewBox="0 0 24 24"><polyline points="22 12 18 12 15 21 9 3 6 12 2 12"/></svg></div></div><div class="stat-value">' + (u.prompt_tokens||0).toLocaleString() + '</div></div>' +
        '<div class="stat-card"><div class="stat-card-header"><div class="stat-label">输出 Tokens</div><div class="stat-card-icon cyan"><svg viewBox="0 0 24 24"><polyline points="22 12 18 12 15 21 9 3 6 12 2 12"/></svg></div></div><div class="stat-value">' + (u.completion_tokens||0).toLocaleString() + '</div></div>' +
        '<div class="stat-card"><div class="stat-card-header"><div class="stat-label">总 Tokens</div><div class="stat-card-icon green"><svg viewBox="0 0 24 24"><circle cx="12" cy="12" r="10"/><path d="M16 8h-6a2 2 0 1 0 0 4h4a2 2 0 1 1 0 4H8"/><path d="M12 18V6"/></svg></div></div><div class="stat-value">' + (u.total_tokens||0).toLocaleString() + '</div></div>';
    }
  } catch(e) {}

  try {
    const r = await fetch(API + '/user/api/usage/details?limit=50', { headers: userHeaders() });
    const d = await r.json();
    if (d.code === 0 && d.data && d.data.length > 0) {
      let html = '<table class="data-table"><thead><tr><th>时间</th><th>模型</th><th>输入</th><th>输出</th><th>总计</th><th>延迟</th><th>状态</th></tr></thead><tbody>';
      d.data.forEach(u => {
        const time = new Date(u.created_at).toLocaleString('zh-CN');
        const status = u.status === 'success' ? '<span class="badge badge-success">成功</span>' : '<span class="badge badge-danger">' + esc(u.status) + '</span>';
        html += '<tr><td>' + esc(time) + '</td><td><strong style="color:var(--text-primary)">' + esc(u.model_name) + '</strong></td><td>' + (u.prompt_tokens||0).toLocaleString() + '</td><td>' + (u.completion_tokens||0).toLocaleString() + '</td><td><strong style="color:var(--text-primary)">' + (u.total_tokens||0).toLocaleString() + '</strong></td><td>' + (u.latency_ms||0) + 'ms</td><td>' + status + '</td></tr>';
      });
      html += '</tbody></table>';
      document.getElementById('usageDetails').innerHTML = html;
    } else {
      document.getElementById('usageDetails').innerHTML = '<div class="empty-state"><div class="empty-state-icon"><svg viewBox="0 0 24 24"><path d="M22 12h-4l-3 9L9 3l-3 9H2"/></svg></div><p>暂无调用记录</p></div>';
    }
  } catch(e) {
    document.getElementById('usageDetails').innerHTML = '<div class="empty-state"><div class="empty-state-icon"><svg viewBox="0 0 24 24"><circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/></svg></div><p>加载失败</p></div>';
  }
}

// ===== Theme =====
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
document.addEventListener('click', function(e) {
  const panel = document.getElementById('themePanel');
  const fab = document.getElementById('themeFab');
  if (panel && fab && !panel.contains(e.target) && !fab.contains(e.target)) {
    panel.classList.remove('open');
    fab.classList.remove('open');
  }
});
function loadTheme() {
  const saved = localStorage.getItem('ui-theme') || 'midnight';
  setTheme(saved);
}

// ===== Mobile Sidebar =====
function openSidebar() {
  document.getElementById('sidebar').classList.add('open');
  document.getElementById('sidebarOverlay').classList.add('show');
  document.body.style.overflow = 'hidden';
}
function closeSidebar() {
  document.getElementById('sidebar').classList.remove('open');
  document.getElementById('sidebarOverlay').classList.remove('show');
  document.body.style.overflow = '';
}
// Auto close sidebar on nav click (mobile)
const _origShowPage = showPage;
showPage = function(name, el) {
  _origShowPage(name, el);
  if (window.innerWidth <= 768) closeSidebar();
};

// ===== Init =====
window.addEventListener('DOMContentLoaded', () => {
  loadTheme();
  if (userToken && currentUser) {
    enterApp();
  } else {
    document.getElementById('authView').style.display = 'flex';
  }
});
</script>
</body>
</html>
` + ""
