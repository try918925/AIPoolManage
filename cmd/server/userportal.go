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
:root{
  --bg-primary:#0a0a0f;
  --bg-secondary:#12121a;
  --bg-card:#1a1a2e;
  --bg-card-hover:#22223a;
  --bg-input:#16162a;
  --border:#2a2a4a;
  --border-focus:#6c5ce7;
  --text-primary:#e8e8f0;
  --text-secondary:#8888aa;
  --text-muted:#555570;
  --accent:#6c5ce7;
  --accent-hover:#7d6ff0;
  --accent-glow:rgba(108,92,231,0.3);
  --success:#00d68f;
  --warning:#f0a030;
  --danger:#ff4757;
  --info:#3498db;
  --gradient-1:linear-gradient(135deg,#6c5ce7,#a855f7);
  --gradient-2:linear-gradient(135deg,#00d68f,#00b894);
  --gradient-3:linear-gradient(135deg,#f0a030,#e17055);
  --gradient-4:linear-gradient(135deg,#3498db,#6c5ce7);
  --shadow-sm:0 2px 8px rgba(0,0,0,0.3);
  --shadow-md:0 4px 16px rgba(0,0,0,0.4);
  --shadow-lg:0 8px 32px rgba(0,0,0,0.5);
  --radius:8px;
  --radius-lg:12px;
  --radius-xl:16px;
  --font:'Inter',system-ui,-apple-system,sans-serif;
}
body{font-family:var(--font);background:var(--bg-primary);color:var(--text-primary);line-height:1.6;overflow-x:hidden}

/* Scrollbar */
::-webkit-scrollbar{width:6px}
::-webkit-scrollbar-track{background:var(--bg-secondary)}
::-webkit-scrollbar-thumb{background:var(--border);border-radius:3px}
::-webkit-scrollbar-thumb:hover{background:var(--accent)}

/* Background animation */
.bg-grid{position:fixed;inset:0;z-index:0;background-image:
  radial-gradient(circle at 20% 50%,rgba(108,92,231,0.08) 0%,transparent 50%),
  radial-gradient(circle at 80% 20%,rgba(0,214,143,0.06) 0%,transparent 50%),
  radial-gradient(circle at 50% 80%,rgba(168,85,247,0.05) 0%,transparent 50%);
  pointer-events:none}

/* ===== Login & Register ===== */
.auth-container{display:flex;align-items:center;justify-content:center;min-height:100vh;position:relative;z-index:1}
.auth-card{background:var(--bg-card);border:1px solid var(--border);border-radius:var(--radius-xl);padding:48px;width:100%;max-width:440px;box-shadow:var(--shadow-lg)}
.auth-card h2{font-size:26px;font-weight:800;text-align:center;margin-bottom:6px;background:var(--gradient-1);-webkit-background-clip:text;-webkit-text-fill-color:transparent}
.auth-card .subtitle{text-align:center;color:var(--text-secondary);margin-bottom:32px;font-size:14px}
.auth-card .form-group{margin-bottom:20px}
.auth-card .form-group label{display:block;font-size:13px;font-weight:600;color:var(--text-secondary);margin-bottom:6px;text-transform:uppercase;letter-spacing:0.5px}
.auth-card .form-input{width:100%;padding:12px 16px;background:var(--bg-input);border:1px solid var(--border);border-radius:var(--radius);color:var(--text-primary);font-size:14px;transition:all .2s}
.auth-card .form-input:focus{outline:none;border-color:var(--accent);box-shadow:0 0 0 3px var(--accent-glow)}
.auth-card .btn-primary{width:100%;padding:13px;font-size:15px;font-weight:700;border:none;border-radius:var(--radius);background:var(--gradient-1);color:#fff;cursor:pointer;transition:all .3s;margin-top:8px;letter-spacing:0.3px}
.auth-card .btn-primary:hover{transform:translateY(-1px);box-shadow:0 4px 20px var(--accent-glow)}
.auth-card .btn-primary:disabled{opacity:0.5;cursor:not-allowed;transform:none}
.auth-card .switch-link{text-align:center;margin-top:24px;font-size:14px;color:var(--text-secondary)}
.auth-card .switch-link a{color:var(--accent);cursor:pointer;text-decoration:none;font-weight:600}
.auth-card .switch-link a:hover{text-decoration:underline}
.auth-card .error-msg{color:var(--danger);font-size:13px;margin-top:8px;text-align:center;min-height:20px}
.auth-card .brand-icon{text-align:center;font-size:48px;margin-bottom:16px}

/* ===== Main App Layout ===== */
.app-container{display:none;min-height:100vh;position:relative;z-index:1}
.app-header{background:var(--bg-secondary);border-bottom:1px solid var(--border);padding:0 32px;height:64px;display:flex;align-items:center;justify-content:space-between;position:sticky;top:0;z-index:100;backdrop-filter:blur(12px)}
.app-header .logo{display:flex;align-items:center;gap:12px;font-weight:800;font-size:18px}
.app-header .logo span{background:var(--gradient-1);-webkit-background-clip:text;-webkit-text-fill-color:transparent}
.app-header .user-info{display:flex;align-items:center;gap:16px}
.app-header .user-badge{display:flex;align-items:center;gap:8px;padding:6px 16px;background:var(--bg-card);border:1px solid var(--border);border-radius:20px;font-size:13px;font-weight:600}
.app-header .user-badge .avatar{width:28px;height:28px;border-radius:50%;background:var(--gradient-1);display:flex;align-items:center;justify-content:center;font-size:13px;font-weight:700;color:#fff}
.btn-logout{padding:6px 16px;background:transparent;border:1px solid var(--border);border-radius:var(--radius);color:var(--text-secondary);cursor:pointer;font-size:13px;transition:all .2s}
.btn-logout:hover{border-color:var(--danger);color:var(--danger)}

/* Nav tabs */
.nav-tabs{display:flex;gap:4px;padding:0 32px;background:var(--bg-secondary);border-bottom:1px solid var(--border)}
.nav-tab{padding:12px 20px;font-size:14px;font-weight:600;color:var(--text-secondary);cursor:pointer;border-bottom:2px solid transparent;transition:all .2s;background:none;border-top:none;border-left:none;border-right:none}
.nav-tab:hover{color:var(--text-primary)}
.nav-tab.active{color:var(--accent);border-bottom-color:var(--accent)}

/* Main content */
.main-content{max-width:1200px;margin:0 auto;padding:32px}
.page{display:none}
.page.active{display:block}
.page-title{font-size:24px;font-weight:800;margin-bottom:8px}
.page-desc{color:var(--text-secondary);font-size:14px;margin-bottom:28px}

/* Stat cards */
.stats-grid{display:grid;grid-template-columns:repeat(auto-fit,minmax(200px,1fr));gap:16px;margin-bottom:32px}
.stat-card{background:var(--bg-card);border:1px solid var(--border);border-radius:var(--radius-lg);padding:24px;transition:all .3s}
.stat-card:hover{border-color:var(--accent);transform:translateY(-2px)}
.stat-card .stat-label{font-size:12px;text-transform:uppercase;letter-spacing:1px;color:var(--text-secondary);margin-bottom:8px;font-weight:600}
.stat-card .stat-value{font-size:28px;font-weight:800}
.stat-card .stat-icon{float:right;font-size:24px;opacity:0.6}

/* Tables */
.data-table{width:100%;border-collapse:collapse;background:var(--bg-card);border:1px solid var(--border);border-radius:var(--radius-lg);overflow:hidden}
.data-table th{background:var(--bg-secondary);padding:14px 16px;text-align:left;font-size:12px;text-transform:uppercase;letter-spacing:1px;color:var(--text-secondary);font-weight:700;border-bottom:1px solid var(--border)}
.data-table td{padding:14px 16px;border-bottom:1px solid var(--border);font-size:13px}
.data-table tr:last-child td{border-bottom:none}
.data-table tr:hover td{background:var(--bg-card-hover)}

/* Status badges */
.badge{display:inline-flex;align-items:center;gap:4px;padding:4px 10px;border-radius:20px;font-size:11px;font-weight:700;text-transform:uppercase;letter-spacing:0.5px}
.badge-success{background:rgba(0,214,143,0.15);color:var(--success)}
.badge-danger{background:rgba(255,71,87,0.15);color:var(--danger)}
.badge-info{background:rgba(52,152,219,0.15);color:var(--info)}
.badge-warning{background:rgba(240,160,48,0.15);color:var(--warning)}

/* Buttons */
.btn{display:inline-flex;align-items:center;gap:6px;padding:8px 16px;border-radius:var(--radius);font-size:13px;font-weight:600;cursor:pointer;transition:all .2s;border:none}
.btn-sm{padding:6px 12px;font-size:12px}
.btn-primary{background:var(--gradient-1);color:#fff}
.btn-primary:hover{transform:translateY(-1px);box-shadow:0 4px 16px var(--accent-glow)}
.btn-danger{background:rgba(255,71,87,0.15);color:var(--danger);border:1px solid rgba(255,71,87,0.3)}
.btn-danger:hover{background:rgba(255,71,87,0.25)}
.btn-ghost{background:transparent;color:var(--text-secondary);border:1px solid var(--border)}
.btn-ghost:hover{border-color:var(--accent);color:var(--accent)}

/* Modal */
.modal-overlay{position:fixed;inset:0;background:rgba(0,0,0,0.6);backdrop-filter:blur(4px);z-index:200;display:none;align-items:center;justify-content:center}
.modal-overlay.show{display:flex}
.modal{background:var(--bg-card);border:1px solid var(--border);border-radius:var(--radius-xl);padding:32px;width:90%;max-width:500px;box-shadow:var(--shadow-lg)}
.modal h3{font-size:20px;font-weight:700;margin-bottom:20px}
.modal .form-group{margin-bottom:16px}
.modal .form-group label{display:block;font-size:12px;font-weight:600;color:var(--text-secondary);margin-bottom:6px;text-transform:uppercase;letter-spacing:.5px}
.modal .form-input{width:100%;padding:10px 14px;background:var(--bg-input);border:1px solid var(--border);border-radius:var(--radius);color:var(--text-primary);font-size:14px;transition:all .2s}
.modal .form-input:focus{outline:none;border-color:var(--accent);box-shadow:0 0 0 3px var(--accent-glow)}
.modal .modal-actions{display:flex;gap:12px;justify-content:flex-end;margin-top:24px}

/* Toast */
.toast{position:fixed;top:20px;right:20px;padding:14px 24px;border-radius:var(--radius);font-size:13px;font-weight:600;z-index:300;transform:translateX(120%);transition:transform .3s;max-width:400px}
.toast.show{transform:translateX(0)}
.toast-success{background:var(--success);color:#000}
.toast-error{background:var(--danger);color:#fff}

/* Key display */
.key-display{background:var(--bg-primary);border:1px solid var(--accent);border-radius:var(--radius);padding:14px 16px;font-family:monospace;font-size:13px;word-break:break-all;margin:12px 0;position:relative}
.key-display .copy-btn{position:absolute;right:8px;top:8px;padding:4px 10px;background:var(--accent);color:#fff;border:none;border-radius:4px;cursor:pointer;font-size:11px;font-weight:700}

/* Code block */
.code-block{background:var(--bg-primary);border:1px solid var(--border);border-radius:var(--radius);padding:16px;font-family:'Consolas','Monaco',monospace;font-size:12px;line-height:1.6;overflow-x:auto;white-space:pre;color:var(--text-secondary);position:relative;margin-top:12px}
.code-block .copy-btn{position:absolute;right:8px;top:8px;padding:4px 10px;background:var(--bg-card);color:var(--text-secondary);border:1px solid var(--border);border-radius:4px;cursor:pointer;font-size:11px}

/* Empty state */
.empty-state{text-align:center;padding:60px 20px;color:var(--text-muted)}
.empty-state .empty-icon{font-size:48px;margin-bottom:16px;opacity:0.5}
.empty-state p{font-size:14px}

/* Progress bar */
.progress-bar{background:var(--bg-primary);border-radius:4px;height:8px;overflow:hidden;margin-top:8px}
.progress-bar .fill{height:100%;border-radius:4px;transition:width .3s}

/* Responsive */
@media(max-width:768px){
  .main-content{padding:20px 16px}
  .stats-grid{grid-template-columns:1fr 1fr}
  .nav-tabs{padding:0 16px;overflow-x:auto}
  .app-header{padding:0 16px}
  .auth-card{margin:16px;padding:32px 24px}
}
</style>
</head>
<body>
<div class="bg-grid"></div>

<!-- Auth Views -->
<div id="authView" class="auth-container">
  <div class="auth-card" id="loginCard">
    <div class="brand-icon">🤖</div>
    <h2>AI Gateway</h2>
    <p class="subtitle">统一 AI 模型 API 网关 · 用户中心</p>
    <div class="form-group">
      <label>用户名</label>
      <input type="text" class="form-input" id="loginUser" placeholder="请输入用户名" autocomplete="username">
    </div>
    <div class="form-group">
      <label>密码</label>
      <input type="password" class="form-input" id="loginPass" placeholder="请输入密码" autocomplete="current-password" onkeypress="if(event.key==='Enter')doLogin()">
    </div>
    <div class="error-msg" id="loginError"></div>
    <button class="btn-primary" onclick="doLogin()" id="loginBtn">🔐 登 录</button>
    <div class="switch-link">还没有账号？<a onclick="showRegister()">立即注册</a></div>
  </div>

  <div class="auth-card" id="registerCard" style="display:none">
    <div class="brand-icon">✨</div>
    <h2>创建账号</h2>
    <p class="subtitle">注册后即可获取专属 API Key</p>
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
    <button class="btn-primary" onclick="doRegister()" id="regBtn">🚀 注 册</button>
    <div class="switch-link">已有账号？<a onclick="showLogin()">返回登录</a></div>
  </div>
</div>

<!-- Main App -->
<div class="app-container" id="appView">
  <header class="app-header">
    <div class="logo">🤖 <span>AI Gateway</span></div>
    <div class="user-info">
      <div class="user-badge">
        <div class="avatar" id="userAvatar">U</div>
        <span id="usernameDisplay">user</span>
      </div>
      <button class="btn-logout" onclick="doLogout()">退出</button>
    </div>
  </header>

  <nav class="nav-tabs">
    <button class="nav-tab active" data-page="dashboard" onclick="showPage('dashboard',this)">📊 总览</button>
    <button class="nav-tab" data-page="apikeys" onclick="showPage('apikeys',this)">🔑 API Keys</button>
    <button class="nav-tab" data-page="models" onclick="showPage('models',this)">🧠 模型列表</button>
    <button class="nav-tab" data-page="usage" onclick="showPage('usage',this)">📈 用量统计</button>
    <button class="nav-tab" data-page="quickstart" onclick="showPage('quickstart',this)">🚀 快速开始</button>
  </nav>

  <div class="main-content">
    <!-- Dashboard -->
    <div class="page active" id="page-dashboard">
      <h1 class="page-title">总览</h1>
      <p class="page-desc">欢迎回来！这里是您的 API 使用概况</p>
      <div class="stats-grid">
        <div class="stat-card">
          <div class="stat-icon">🔑</div>
          <div class="stat-label">API Keys</div>
          <div class="stat-value" id="statKeys">-</div>
        </div>
        <div class="stat-card">
          <div class="stat-icon">📊</div>
          <div class="stat-label">总调用次数</div>
          <div class="stat-value" id="statCalls">-</div>
        </div>
        <div class="stat-card">
          <div class="stat-icon">🪙</div>
          <div class="stat-label">总 Token 消耗</div>
          <div class="stat-value" id="statTokens">-</div>
        </div>
        <div class="stat-card">
          <div class="stat-icon">🧠</div>
          <div class="stat-label">可用模型数</div>
          <div class="stat-value" id="statModels">-</div>
        </div>
      </div>
      <h3 style="margin-bottom:16px;font-weight:700">最近调用</h3>
      <div id="recentUsage"></div>
    </div>

    <!-- API Keys -->
    <div class="page" id="page-apikeys">
      <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:24px">
        <div>
          <h1 class="page-title" style="margin-bottom:4px">API Keys</h1>
          <p class="page-desc" style="margin-bottom:0">管理您的 API 访问密钥</p>
        </div>
        <button class="btn btn-primary" onclick="showCreateKeyModal()">+ 创建 Key</button>
      </div>
      <div id="keysList"></div>
    </div>

    <!-- Models -->
    <div class="page" id="page-models">
      <h1 class="page-title">可用模型</h1>
      <p class="page-desc">以下是当前网关支持的所有 AI 模型</p>
      <div id="modelsList"></div>
    </div>

    <!-- Usage -->
    <div class="page" id="page-usage">
      <h1 class="page-title">用量统计</h1>
      <p class="page-desc">查看您的 API 调用详情和 Token 消耗</p>
      <div class="stats-grid" id="usageStats"></div>
      <h3 style="margin-bottom:16px;font-weight:700">调用明细</h3>
      <div id="usageDetails"></div>
    </div>

    <!-- Quick Start -->
    <div class="page" id="page-quickstart">
      <h1 class="page-title">快速开始</h1>
      <p class="page-desc">几分钟内开始使用 AI Gateway API</p>

      <div style="background:var(--bg-card);border:1px solid var(--border);border-radius:var(--radius-lg);padding:28px;margin-bottom:20px">
        <h3 style="margin-bottom:12px;font-weight:700">📖 第一步：获取 API Key</h3>
        <p style="color:var(--text-secondary);font-size:14px;margin-bottom:12px">在 <a style="color:var(--accent);cursor:pointer" onclick="showPage('apikeys',document.querySelector('[data-page=apikeys]'))">API Keys</a> 页面创建一个新的密钥</p>
      </div>

      <div style="background:var(--bg-card);border:1px solid var(--border);border-radius:var(--radius-lg);padding:28px;margin-bottom:20px">
        <h3 style="margin-bottom:12px;font-weight:700">🔗 第二步：调用 API</h3>
        <p style="color:var(--text-secondary);font-size:14px;margin-bottom:4px">API 兼容 OpenAI 格式，可直接使用 OpenAI SDK</p>
        <p style="color:var(--text-secondary);font-size:14px;margin-bottom:12px">Base URL: <code style="background:var(--bg-primary);padding:2px 8px;border-radius:4px;color:var(--accent)">${location.origin}/v1</code></p>
        <h4 style="margin-bottom:8px;font-size:13px;color:var(--text-secondary);text-transform:uppercase;letter-spacing:1px">cURL 示例</h4>
        <div class="code-block" id="curlExample"><button class="copy-btn" onclick="copyCode('curlExample')">复制</button>curl ${location.origin}/v1/chat/completions \\
  -H "Authorization: Bearer YOUR_API_KEY" \\
  -H "Content-Type: application/json" \\
  -d '{
  "model": "gpt-4o",
  "messages": [{"role": "user", "content": "Hello!"}]
}'</div>
        <h4 style="margin:20px 0 8px;font-size:13px;color:var(--text-secondary);text-transform:uppercase;letter-spacing:1px">Python (OpenAI SDK)</h4>
        <div class="code-block" id="pythonExample"><button class="copy-btn" onclick="copyCode('pythonExample')">复制</button>from openai import OpenAI

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

      <div style="background:var(--bg-card);border:1px solid var(--border);border-radius:var(--radius-lg);padding:28px">
        <h3 style="margin-bottom:12px;font-weight:700">📋 第三步：查看模型</h3>
        <p style="color:var(--text-secondary);font-size:14px">在 <a style="color:var(--accent);cursor:pointer" onclick="showPage('models',document.querySelector('[data-page=models]'))">模型列表</a> 查看所有可用模型及其详情</p>
      </div>
    </div>
  </div>
</div>

<!-- Create Key Modal -->
<div class="modal-overlay" id="createKeyModal">
  <div class="modal">
    <h3>🔑 创建 API Key</h3>
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
    <h3>✅ Key 创建成功</h3>
    <p style="color:var(--danger);font-size:13px;margin-bottom:8px">⚠️ 请立即复制保存，此密钥仅显示一次！</p>
    <div class="key-display" id="newKeyDisplay">
      <button class="copy-btn" onclick="copyNewKey()">复制</button>
      <span id="newKeyValue"></span>
    </div>
    <div class="modal-actions">
      <button class="btn btn-primary" onclick="closeModal('keyCreatedModal');loadKeys()">我已保存</button>
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
  t.textContent = msg;
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
    btn.textContent = '🔐 登 录';
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
    btn.textContent = '🚀 注 册';
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
  document.querySelectorAll('.nav-tab').forEach(t => t.classList.remove('active'));
  document.getElementById('page-'+name).classList.add('active');
  if (el) el.classList.add('active');

  if (name === 'dashboard') loadDashboard();
  else if (name === 'apikeys') loadKeys();
  else if (name === 'models') loadModels();
  else if (name === 'usage') loadUsage();
}

// ===== Dashboard =====
async function loadDashboard() {
  // Load keys count
  try {
    const r = await fetch(API + '/user/api/keys', { headers: userHeaders() });
    const d = await r.json();
    if (d.code === 0) {
      const keys = d.data || [];
      document.getElementById('statKeys').textContent = keys.length;
    }
  } catch(e) {}

  // Load usage summary
  try {
    const r = await fetch(API + '/user/api/usage', { headers: userHeaders() });
    const d = await r.json();
    if (d.code === 0 && d.data) {
      document.getElementById('statCalls').textContent = (d.data.total_calls || 0).toLocaleString();
      document.getElementById('statTokens').textContent = (d.data.total_tokens || 0).toLocaleString();
    }
  } catch(e) {}

  // Load models count
  try {
    const r = await fetch(API + '/user/api/models', { headers: userHeaders() });
    const d = await r.json();
    const mlist = d.data || [];
    document.getElementById('statModels').textContent = mlist.length || 0;
  } catch(e) {}

  // Load recent usage
  try {
    const r = await fetch(API + '/user/api/usage/details?limit=5', { headers: userHeaders() });
    const d = await r.json();
    if (d.code === 0 && d.data && d.data.length > 0) {
      let html = '<table class="data-table"><thead><tr><th>时间</th><th>模型</th><th>Tokens</th><th>延迟</th><th>状态</th></tr></thead><tbody>';
      d.data.forEach(u => {
        const time = new Date(u.created_at).toLocaleString('zh-CN');
        const status = u.status === 'success' ? '<span class="badge badge-success">成功</span>' : '<span class="badge badge-danger">失败</span>';
        html += '<tr><td>' + esc(time) + '</td><td>' + esc(u.model_name) + '</td><td>' + (u.total_tokens||0).toLocaleString() + '</td><td>' + (u.latency_ms||0) + 'ms</td><td>' + status + '</td></tr>';
      });
      html += '</tbody></table>';
      document.getElementById('recentUsage').innerHTML = html;
    } else {
      document.getElementById('recentUsage').innerHTML = '<div class="empty-state"><div class="empty-icon">📭</div><p>暂无调用记录<br>创建 API Key 后开始使用吧</p></div>';
    }
  } catch(e) {
    document.getElementById('recentUsage').innerHTML = '<div class="empty-state"><div class="empty-icon">📭</div><p>暂无调用记录</p></div>';
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
        document.getElementById('keysList').innerHTML = '<div class="empty-state"><div class="empty-icon">🔑</div><p>您还没有 API Key<br>点击上方按钮创建第一个</p></div>';
        return;
      }
      let html = '<table class="data-table"><thead><tr><th>名称</th><th>Key 前缀</th><th>速率限制</th><th>配额</th><th>状态</th><th>创建时间</th><th>操作</th></tr></thead><tbody>';
      keys.forEach(k => {
        const enabled = k.enabled ? '<span class="badge badge-success">启用</span>' : '<span class="badge badge-danger">禁用</span>';
        const quota = k.quota_limit > 0 ? ((k.quota_used||0).toLocaleString() + ' / ' + k.quota_limit.toLocaleString()) : '无限';
        const time = new Date(k.created_at).toLocaleDateString('zh-CN');
        html += '<tr><td><strong>' + esc(k.name||'未命名') + '</strong></td><td><code style="background:var(--bg-primary);padding:2px 6px;border-radius:4px">' + esc(k.key_prefix) + '...</code></td><td>' + k.rate_limit + '/min</td><td>' + quota + '</td><td>' + enabled + '</td><td>' + esc(time) + '</td><td><button class="btn btn-danger btn-sm" onclick="deleteKey(' + k.id + ')">删除</button></td></tr>';
      });
      html += '</tbody></table>';
      document.getElementById('keysList').innerHTML = html;
    }
  } catch(e) {
    document.getElementById('keysList').innerHTML = '<div class="empty-state"><div class="empty-icon">⚠️</div><p>加载失败，请刷新重试</p></div>';
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
      let html = '<div style="display:grid;grid-template-columns:repeat(auto-fill,minmax(300px,1fr));gap:16px">';
      models.forEach(m => {
        const name = m.id || m.model_name || m.name || 'unknown';
        const provider = m.owned_by || m.provider_type || 'unknown';
        const mtype = m.model_type || 'chat';
        const channels = m.channel_count || 1;
        html += '<div style="background:var(--bg-card);border:1px solid var(--border);border-radius:var(--radius-lg);padding:24px;transition:all .3s" onmouseenter="this.style.borderColor=\'var(--accent)\'" onmouseleave="this.style.borderColor=\'var(--border)\'">' +
          '<div style="display:flex;justify-content:space-between;align-items:start;margin-bottom:12px">' +
          '<h3 style="font-size:16px;font-weight:700">' + esc(name) + '</h3>' +
          '<span class="badge badge-info">' + esc(provider) + '</span>' +
          '</div>' +
          '<div style="color:var(--text-secondary);font-size:13px;display:flex;gap:16px">' +
          '<p>类型: <span style="color:var(--text-primary)">' + esc(mtype) + '</span></p>' +
          '<p>通道: <span style="color:var(--success)">' + channels + '</span></p>' +
          '</div></div>';
      });
      html += '</div>';
      document.getElementById('modelsList').innerHTML = html;
    } else {
      document.getElementById('modelsList').innerHTML = '<div class="empty-state"><div class="empty-icon">🧠</div><p>暂无可用模型<br>请联系管理员配置</p></div>';
    }
  } catch(e) {
    document.getElementById('modelsList').innerHTML = '<div class="empty-state"><div class="empty-icon">⚠️</div><p>加载失败</p></div>';
  }
}

// ===== Usage =====
async function loadUsage() {
  // Summary
  try {
    const r = await fetch(API + '/user/api/usage', { headers: userHeaders() });
    const d = await r.json();
    if (d.code === 0 && d.data) {
      const u = d.data;
      document.getElementById('usageStats').innerHTML =
        '<div class="stat-card"><div class="stat-icon">📊</div><div class="stat-label">总调用</div><div class="stat-value">' + (u.total_calls||0).toLocaleString() + '</div></div>' +
        '<div class="stat-card"><div class="stat-icon">📥</div><div class="stat-label">输入 Tokens</div><div class="stat-value">' + (u.prompt_tokens||0).toLocaleString() + '</div></div>' +
        '<div class="stat-card"><div class="stat-icon">📤</div><div class="stat-label">输出 Tokens</div><div class="stat-value">' + (u.completion_tokens||0).toLocaleString() + '</div></div>' +
        '<div class="stat-card"><div class="stat-icon">🪙</div><div class="stat-label">总 Tokens</div><div class="stat-value">' + (u.total_tokens||0).toLocaleString() + '</div></div>';
    }
  } catch(e) {}

  // Details
  try {
    const r = await fetch(API + '/user/api/usage/details?limit=50', { headers: userHeaders() });
    const d = await r.json();
    if (d.code === 0 && d.data && d.data.length > 0) {
      let html = '<table class="data-table"><thead><tr><th>时间</th><th>模型</th><th>输入</th><th>输出</th><th>总计</th><th>延迟</th><th>状态</th></tr></thead><tbody>';
      d.data.forEach(u => {
        const time = new Date(u.created_at).toLocaleString('zh-CN');
        const status = u.status === 'success' ? '<span class="badge badge-success">成功</span>' : '<span class="badge badge-danger">' + esc(u.status) + '</span>';
        html += '<tr><td>' + esc(time) + '</td><td>' + esc(u.model_name) + '</td><td>' + (u.prompt_tokens||0).toLocaleString() + '</td><td>' + (u.completion_tokens||0).toLocaleString() + '</td><td><strong>' + (u.total_tokens||0).toLocaleString() + '</strong></td><td>' + (u.latency_ms||0) + 'ms</td><td>' + status + '</td></tr>';
      });
      html += '</tbody></table>';
      document.getElementById('usageDetails').innerHTML = html;
    } else {
      document.getElementById('usageDetails').innerHTML = '<div class="empty-state"><div class="empty-icon">📭</div><p>暂无调用记录</p></div>';
    }
  } catch(e) {
    document.getElementById('usageDetails').innerHTML = '<div class="empty-state"><div class="empty-icon">⚠️</div><p>加载失败</p></div>';
  }
}

// ===== Init =====
window.addEventListener('DOMContentLoaded', () => {
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
