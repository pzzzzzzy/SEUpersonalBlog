// 全局变量
let currentUser = null;
let currentArticleId = null;
let authToken = localStorage.getItem('authToken');
// API基础URL配置 - 指向Traefik代理
const API_BASE_URL = 'http://localhost:8086'; // Traefik代理，用于访问后端API

// 示例文章数据
const sampleArticles = [
    {
        id: 1,
        title: "欢迎来到我的博客",
        content: "这是我的第一篇博客文章！在这里我会分享我的技术心得、生活感悟和学习体会。希望能和大家一起成长，共同进步。",
        summary: "博客开篇，分享技术心得和生活感悟",
        author: { username: "博主小明", id: 1 },
        view_count: 128,
        comment_count: 5,
        created_at: "2024-01-15T10:30:00Z",
        tags: ["生活", "技术"],
        liked: false,
        likes: 23
    },
    {
        id: 2,
        title: "Go语言学习笔记",
        content: "最近在学习Go语言，发现它的并发模型真的很强大。goroutine和channel的设计让并发编程变得简单易懂。分享一下我的学习心得...",
        summary: "Go语言并发编程学习心得",
        author: { username: "技术达人", id: 2 },
        view_count: 256,
        comment_count: 8,
        created_at: "2024-01-20T14:20:00Z",
        tags: ["Go语言", "编程", "并发"],
        liked: false,
        likes: 45
    },
    {
        id: 3,
        title: "前端开发最佳实践",
        content: "在前端开发中，有很多最佳实践值得我们遵循。比如组件化开发、状态管理、性能优化等。今天想和大家分享一些实用的技巧...",
        summary: "前端开发的最佳实践和实用技巧",
        author: { username: "前端小王", id: 3 },
        view_count: 189,
        comment_count: 12,
        created_at: "2024-01-25T09:15:00Z",
        tags: ["前端", "JavaScript", "最佳实践"],
        liked: true,
        likes: 67
    }
];

// 示例评论数据
const sampleComments = {
    1: [
        {
            id: 1,
            content: "写得很好，期待更多分享！",
            user: { username: "读者A" },
            created_at: "2024-01-15T11:00:00Z",
            likes: 5,
            liked: false
        },
        {
            id: 2,
            content: "博客界面很简洁，喜欢这种风格。",
            user: { username: "设计爱好者" },
            created_at: "2024-01-15T12:30:00Z",
            likes: 3,
            liked: false
        }
    ],
    2: [
        {
            id: 3,
            content: "Go语言确实很适合并发编程！",
            user: { username: "后端开发者" },
            created_at: "2024-01-20T15:00:00Z",
            likes: 8,
            liked: true
        }
    ],
    3: [
        {
            id: 4,
            content: "这些技巧很实用，已经用上了！",
            user: { username: "前端新手" },
            created_at: "2024-01-25T10:30:00Z",
            likes: 12,
            liked: false
        }
    ]
};

// 页面切换函数
function showPage(pageId) {
    // 隐藏所有页面
    document.querySelectorAll('.page').forEach(page => {
        page.style.display = 'none';
    });
    
    // 显示指定页面
    document.getElementById(pageId).style.display = 'block';
}

// 显示首页
function showHome() {
    showPage('home-page');
    loadRecentArticles();
}

// 显示文章列表
function showArticles() {
    showPage('articles-page');
    loadArticles();
}

// 显示登录页
function showLogin() {
    showPage('login-page');
}

// 显示注册页
function showRegister() {
    showPage('register-page');
}

// 显示个人中心
function showProfile() {
    if (!currentUser) {
        showLogin();
        return;
    }
    showPage('profile-page');
    loadProfile();
}

// 显示写文章页面
function showCreateArticle() {
    if (!currentUser) {
        showLogin();
        return;
    }
    showPage('create-article-page');
    setupContentChecker(); // 启动内容检查器
}

// 显示我的文章
function showMyArticles() {
    if (!currentUser) {
        showLogin();
        return;
    }
    showPage('articles-page');
    loadMyArticles();
}

// 初始化应用
document.addEventListener('DOMContentLoaded', function() {
    checkAuthStatus();
    setupEventListeners();
    showHome();
});

// 检查认证状态
async function checkAuthStatus() {
    if (authToken) {
        try {
            // 这里可以添加验证token的API调用
            updateUIForLoggedInUser();
        } catch (error) {
            localStorage.removeItem('authToken');
            authToken = null;
        }
    }
}

// 更新UI为登录状态
function updateUIForLoggedInUser() {
    document.getElementById('user-menu').style.display = 'block';
    document.getElementById('guest-menu').style.display = 'none';
    document.getElementById('comment-form').style.display = 'block';
}

// 更新UI为登出状态
function updateUIForGuestUser() {
    document.getElementById('user-menu').style.display = 'none';
    document.getElementById('guest-menu').style.display = 'block';
    document.getElementById('comment-form').style.display = 'none';
}

// 实时检查文章内容
function setupContentChecker() {
    const contentElement = document.getElementById('article-content');
    const debugDiv = document.getElementById('content-debug');
    
    if (contentElement && debugDiv) {
        contentElement.addEventListener('input', function() {
            const content = this.value;
            debugDiv.style.display = 'block';
            debugDiv.innerHTML = `
                <strong>内容检查：</strong><br>
                长度：${content.length}<br>
                是否为空：${content.length === 0}<br>
                原始值：${JSON.stringify(content.substring(0, 50))}${content.length > 50 ? '...' : ''}
            `;
        });
    }
}

// 设置事件监听器
function setupEventListeners() {
    // 登录表单
    const loginForm = document.getElementById('login-form');
    if (loginForm) {
        loginForm.addEventListener('submit', async function(e) {
            e.preventDefault();
            await login();
        });
    }

    // 注册表单
    const registerForm = document.getElementById('register-form');
    if (registerForm) {
        registerForm.addEventListener('submit', async function(e) {
            e.preventDefault();
            await register();
        });
    }

    // 文章表单
    const articleForm = document.getElementById('article-form');
    if (articleForm) {
        articleForm.addEventListener('submit', async function(e) {
            e.preventDefault();
            await createArticle();
        });
    }
}

// 用户登录
async function login() {
    const username = document.getElementById('login-username').value;
    const password = document.getElementById('login-password').value;

    try {
        const response = await fetch(`${API_BASE_URL}/api/auth/login`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ username, password }),
        });

        if (response.ok) {
            const data = await response.json();
            authToken = data.token;
            currentUser = data.user;
            localStorage.setItem('authToken', authToken);
            
            updateUIForLoggedInUser();
            showHome();
            alert('登录成功！');
        } else {
            const error = await response.json();
            alert('登录失败：' + error.error);
        }
    } catch (error) {
        alert('登录失败：网络错误');
    }
}

// 用户注册
async function register() {
    const username = document.getElementById('register-username').value;
    const email = document.getElementById('register-email').value;
    const password = document.getElementById('register-password').value;

    try {
        const response = await fetch(`${API_BASE_URL}/api/auth/register`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ username, email, password }),
        });

        if (response.ok) {
            const data = await response.json();
            authToken = data.token;
            currentUser = data.user;
            localStorage.setItem('authToken', authToken);
            
            updateUIForLoggedInUser();
            showHome();
            alert('注册成功！');
        } else {
            const error = await response.json();
            alert('注册失败：' + error.error);
        }
    } catch (error) {
        alert('注册失败：网络错误');
    }
}

// 用户登出
function logout() {
    authToken = null;
    currentUser = null;
    localStorage.removeItem('authToken');
    updateUIForGuestUser();
    showHome();
    alert('已登出');
}

// 加载最新文章
async function loadRecentArticles() {
    try {
        // 暂时直接使用示例数据，避免API调用问题
        console.log('直接使用示例文章数据');
        displayArticles(sampleArticles, 'recent-articles-list');
        
        // 注释掉API调用部分，稍后调试
        /*
        const response = await fetch(`${API_BASE_URL}/api/articles`);
        console.log('API响应状态:', response.status);
        if (response.ok) {
            const data = await response.json();
            console.log('API响应数据:', data);
            
            // 确保articles是数组格式
            let articles = [];
            if (Array.isArray(data)) {
                articles = data;
            } else if (data && data.articles && Array.isArray(data.articles)) {
                articles = data.articles;
            } else {
                console.warn('API返回的数据格式不符合预期，使用示例数据');
                articles = sampleArticles;
            }
            
            displayArticles(articles, 'recent-articles-list');
        } else {
            console.error('加载最新文章失败，状态码:', response.status);
            // 如果API调用失败，回退到示例数据
            displayArticles(sampleArticles, 'recent-articles-list');
        }
        */
    } catch (error) {
        console.error('加载最新文章错误：', error);
        // 如果发生错误，回退到示例数据
        displayArticles(sampleArticles, 'recent-articles-list');
    }
}

// 加载文章列表
async function loadArticles() {
    try {
        const response = await fetch(`${API_BASE_URL}/api/articles`);
        if (response.ok) {
            const articles = await response.json();
            displayArticles(articles, 'articles-list', true);
        } else {
            console.error('加载文章列表失败');
            // 如果API调用失败，回退到示例数据
            displayArticles(sampleArticles, 'articles-list', true);
        }
    } catch (error) {
        console.error('加载文章列表错误：', error);
        // 如果发生错误，回退到示例数据
        displayArticles(sampleArticles, 'articles-list', true);
    }
}

// 加载我的文章
async function loadMyArticles() {
    if (!currentUser) return;
    
    try {
        const response = await fetch(`${API_BASE_URL}/api/articles?author_id=${currentUser.id}`);
        if (response.ok) {
            const articles = await response.json();
            displayArticles(articles, 'articles-list', true);
        }
    } catch (error) {
        console.error('加载文章失败：', error);
    }
}

// 显示文章列表
function displayArticles(articles, containerId, isList = false) {
    const container = document.getElementById(containerId);
    container.innerHTML = '';

    // 确保articles是数组
    if (!Array.isArray(articles)) {
        console.error('传入的articles参数不是数组:', articles);
        container.innerHTML = '<p>加载文章时发生错误</p>';
        return;
    }

    if (articles.length === 0) {
        container.innerHTML = '<p>暂无文章</p>';
        return;
    }

    articles.forEach(article => {
        const articleElement = document.createElement('div');
        articleElement.className = isList ? 'article-item' : 'article-card';
        
        articleElement.innerHTML = `
            <div class="article-header">
                <div class="article-author">
                    <img id="author-avatar-${article.id}" src="/static/images/default-avatar.png" alt="头像" class="author-avatar" style="width: 10%">
                    <span class="author-name">${article.author.username}</span>
                </div>
                <div class="article-actions">
                    <button class="like-btn ${article.liked ? 'liked' : ''}" onclick="toggleLike(${article.id})">
                        <span class="heart">${article.liked ? '❤️' : '🤍'}</span>
                        <span class="like-count">${article.like_count || article.likes || 0}</span>
                    </button>
                </div>
            </div>
            <h4><a href="#" onclick="showArticleDetail(${article.id})">${article.title}</a></h4>
            <p>${article.summary || article.content.substring(0, 100) + '...'}</p>
            <div class="article-tags">
                ${article.tags.map(tag => `<span class="tag">${typeof tag === 'object' ? tag.name : tag}</span>`).join('')}
            </div>
            <div class="article-meta">
                <span>📖 ${article.view_count}</span>
                <span>💬 ${article.comment_count}</span>
                <span>📅 ${new Date(article.created_at).toLocaleDateString()}</span>
            </div>
            <div class="article-comments" id="comments-${article.id}">
                ${renderComments(article.id)}
            </div>
            <div class="comment-form" id="comment-form-${article.id}">
                <textarea placeholder="写下你的评论..." id="comment-input-${article.id}"></textarea>
                <button onclick="addComment(${article.id})">发表评论</button>
            </div>
        `;
        
        container.appendChild(articleElement);
    });
}

// 渲染评论
function renderComments(articleId) {
    const comments = sampleComments[articleId] || [];
    if (comments.length === 0) {
        return '<p style="color: #999; font-size: 14px;">暂无评论，快来发表第一条评论吧！</p>';
    }
    
    return comments.map(comment => `
        <div class="comment-item">
            <div class="comment-header">
                <span class="comment-author">${comment.user.username}</span>
                <span class="comment-time">${new Date(comment.created_at).toLocaleString()}</span>
            </div>
            <div class="comment-content">${comment.content}</div>
            <div class="comment-actions">
                <button class="comment-like-btn ${comment.liked ? 'liked' : ''}" onclick="toggleCommentLike(${articleId}, ${comment.id})">
                    <span class="heart">${comment.liked ? '❤️' : '🤍'}</span>
                    <span>${comment.likes}</span>
                </button>
            </div>
        </div>
    `).join('');
}

// 切换文章点赞状态
async function toggleLike(articleId) {
    if (!currentUser) {
        alert('请先登录后再点赞');
        return;
    }
    
    // 立即更新UI状态，提供即时反馈
    const likeButtons = document.querySelectorAll(`button[onclick="toggleLike(${articleId})"]`);
    let newLikedState = false;
    let newLikeCount = 0;
    
    // 保存原始状态，用于出错时恢复
    const originalStates = [];
    
    likeButtons.forEach(button => {
        const heartSpan = button.querySelector('.heart');
        const likeCountSpan = button.querySelector('.like-count');
        
        // 保存原始状态
        originalStates.push({
            button,
            wasLiked: button.classList.contains('liked'),
            originalCount: parseInt(likeCountSpan?.textContent || '0')
        });
        
        // 切换当前状态（暂时的前端显示）
        if (button.classList.contains('liked')) {
            button.classList.remove('liked');
            if (heartSpan) heartSpan.textContent = '🤍';
            newLikedState = false;
            newLikeCount = parseInt(likeCountSpan?.textContent || '0') - 1;
        } else {
            button.classList.add('liked');
            if (heartSpan) heartSpan.textContent = '❤️';
            newLikedState = true;
            newLikeCount = parseInt(likeCountSpan?.textContent || '0') + 1;
        }
        
        // 立即更新点赞数显示
        if (likeCountSpan) {
            likeCountSpan.textContent = newLikeCount;
        }
    });
    
    // 立即更新文章详情页（如果当前正在查看该文章）
    let detailLikeElements = [];
    if (currentArticleId === articleId) {
        const likeCountElement = document.querySelector('#article-content .like-count');
        const articleMetaLikeCount = document.querySelector('#article-content .article-meta span:nth-child(5)');
        const detailLikeButton = document.querySelector('#article-content button[onclick="toggleLike(' + articleId + ')"]');
        
        detailLikeElements = [likeCountElement, articleMetaLikeCount, detailLikeButton];
        
        if (likeCountElement) {
            likeCountElement.dataset.originalValue = likeCountElement.textContent;
            likeCountElement.textContent = newLikeCount;
        }
        
        if (articleMetaLikeCount) {
            articleMetaLikeCount.dataset.originalValue = articleMetaLikeCount.textContent;
            articleMetaLikeCount.textContent = `点赞：${newLikeCount}`;
        }
        
        if (detailLikeButton) {
            detailLikeButton.dataset.wasLiked = detailLikeButton.classList.contains('liked');
            const heartSpan = detailLikeButton.querySelector('.heart');
            if (heartSpan) {
                heartSpan.dataset.originalText = heartSpan.textContent;
                heartSpan.textContent = newLikedState ? '❤️' : '🤍';
            }
            detailLikeButton.classList.toggle('liked', newLikedState);
        }
    }
    
    try {
        const response = await fetch(`${API_BASE_URL}/api/articles/${articleId}/like`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': authToken ? `Bearer ${authToken}` : ''
            }
        });
        
        if (response.ok) {
            const data = await response.json();
            
            // 使用后端返回的最终状态更新UI
            likeButtons.forEach(button => {
                const heartSpan = button.querySelector('.heart');
                const likeCountSpan = button.querySelector('.like-count');
                
                // 根据后端返回的liked状态设置样式和图标
                if (data.liked) {
                    button.classList.add('liked');
                    if (heartSpan) heartSpan.textContent = '❤️';
                } else {
                    button.classList.remove('liked');
                    if (heartSpan) heartSpan.textContent = '🤍';
                }
                
                // 更新点赞数为后端返回的准确值
                if (likeCountSpan) {
                    likeCountSpan.textContent = data.like_count || 0;
                }
            });
            
            // 只在文章详情页更新特定内容，不重新加载整个页面
            if (currentArticleId === articleId) {
                const likeCountElement = document.querySelector('#article-content .like-count');
                if (likeCountElement) {
                    likeCountElement.textContent = data.like_count || 0;
                    delete likeCountElement.dataset.originalValue;
                }
                
                const articleMetaLikeCount = document.querySelector('#article-content .article-meta span:nth-child(5)');
                if (articleMetaLikeCount) {
                    articleMetaLikeCount.textContent = `点赞：${data.like_count || 0}`;
                    delete articleMetaLikeCount.dataset.originalValue;
                }
                
                // 更新文章详情页的点赞按钮
                const detailLikeButton = document.querySelector('#article-content button[onclick="toggleLike(' + articleId + ')"]');
                if (detailLikeButton) {
                    const heartSpan = detailLikeButton.querySelector('.heart');
                    if (heartSpan) {
                        heartSpan.textContent = data.liked ? '❤️' : '🤍';
                        delete heartSpan.dataset.originalText;
                    }
                    detailLikeButton.classList.toggle('liked', data.liked);
                    delete detailLikeButton.dataset.wasLiked;
                }
            }
        } else {
            const errorData = await response.json().catch(() => ({}));
            throw new Error(errorData.error || '操作失败');
        }
    } catch (error) {
        console.error('点赞操作错误：', error);
        
        // 出错时恢复UI状态
        originalStates.forEach(({button, wasLiked, originalCount}) => {
            const heartSpan = button.querySelector('.heart');
            const likeCountSpan = button.querySelector('.like-count');
            
            // 恢复原来的状态
            if (wasLiked) {
                button.classList.add('liked');
                if (heartSpan) heartSpan.textContent = '❤️';
            } else {
                button.classList.remove('liked');
                if (heartSpan) heartSpan.textContent = '🤍';
            }
            
            // 恢复原来的点赞数
            if (likeCountSpan) {
                likeCountSpan.textContent = originalCount;
            }
        });
        
        // 恢复文章详情页
        if (currentArticleId === articleId) {
            const likeCountElement = document.querySelector('#article-content .like-count');
            if (likeCountElement && likeCountElement.dataset.originalValue) {
                likeCountElement.textContent = likeCountElement.dataset.originalValue;
                delete likeCountElement.dataset.originalValue;
            }
            
            const articleMetaLikeCount = document.querySelector('#article-content .article-meta span:nth-child(5)');
            if (articleMetaLikeCount && articleMetaLikeCount.dataset.originalValue) {
                articleMetaLikeCount.textContent = articleMetaLikeCount.dataset.originalValue;
                delete articleMetaLikeCount.dataset.originalValue;
            }
            
            const detailLikeButton = document.querySelector('#article-content button[onclick="toggleLike(' + articleId + ')"]');
            if (detailLikeButton) {
                const heartSpan = detailLikeButton.querySelector('.heart');
                if (heartSpan && heartSpan.dataset.originalText) {
                    heartSpan.textContent = heartSpan.dataset.originalText;
                    delete heartSpan.dataset.originalText;
                }
                detailLikeButton.classList.toggle('liked', detailLikeButton.dataset.wasLiked === 'true');
                delete detailLikeButton.dataset.wasLiked;
            }
        }
        
        alert('操作出错，请稍后重试');
    }
}

// 切换评论点赞状态
async function toggleCommentLike(articleId, commentId) {
    if (!currentUser) {
        alert('请先登录后再点赞');
        return;
    }
    
    try {
        const response = await fetch(`${API_BASE_URL}/api/comments/${commentId}/like`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': authToken ? `Bearer ${authToken}` : ''
            }
        });
        
        if (response.ok) {
            // 重新加载评论以更新点赞状态
            await loadComments(articleId);
        } else {
            alert('操作失败，请稍后重试');
        }
    } catch (error) {
        console.error('评论点赞操作错误：', error);
        alert('操作出错，请稍后重试');
    }
}

// 添加评论
async function addComment(articleId) {
    if (!currentUser) {
        alert('请先登录后再评论');
        return;
    }
    
    const commentInput = document.getElementById(`comment-input-${articleId}`);
    const content = commentInput.value.trim();
    
    if (!content) {
        alert('请输入评论内容');
        return;
    }
    
    if (content.length < 2) {
        alert('评论内容太短，请至少输入2个字符');
        return;
    }
    
    try {
        const response = await fetch(`${API_BASE_URL}/api/comments`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': authToken ? `Bearer ${authToken}` : ''
            },
            body: JSON.stringify({
                article_id: articleId,
                content: content
            })
        });
        
        if (response.ok) {
            // 清空输入框
            commentInput.value = '';
            // 重新加载评论
            await loadComments(articleId);
            alert('评论发表成功！');
        } else {
            alert('评论发表失败，请稍后重试');
        }
    } catch (error) {
        console.error('发表评论错误：', error);
        alert('发表评论出错，请稍后重试');
    }
    return; // 提前返回，不再执行旧的示例代码
    
    // 下面是旧的示例代码（不会执行）
    const newComment = {
        id: Date.now(), // 使用时间戳作为ID
        content: content,
        user: { username: currentUser.username },
        created_at: new Date().toISOString(),
        likes: 0,
        liked: false
    };
    
    // 添加到示例数据
    if (!sampleComments[articleId]) {
        sampleComments[articleId] = [];
    }
    sampleComments[articleId].unshift(newComment); // 添加到最前面
    
    // 更新评论数
    const article = sampleArticles.find(a => a.id === articleId);
    if (article) {
        article.comment_count++;
    }
    
    // 清空输入框
    commentInput.value = '';
    
    // 重新渲染评论
    const commentsContainer = document.getElementById(`comments-${articleId}`);
    if (commentsContainer) {
        commentsContainer.innerHTML = renderComments(articleId);
    }
    
    // 重新渲染文章（更新评论数）
    displayArticles(sampleArticles, 'recent-articles-list');
    
    alert('评论发表成功！');
}

// 显示文章详情
async function showArticleDetail(articleId) {
    currentArticleId = articleId;
    showPage('article-detail-page');
    
    try {
        const response = await fetch(`${API_BASE_URL}/api/articles/${articleId}`);
        if (response.ok) {
            const article = await response.json();
            displayArticleDetail(article);
            loadComments(articleId);
        }
    } catch (error) {
        console.error('加载文章详情失败：', error);
    }
}

// 显示文章详情内容
function displayArticleDetail(article) {
    const contentElement = document.getElementById('article-content');
    contentElement.innerHTML = `
        <h1>${article.title}</h1>
        <div class="article-meta">
            <span>作者：${article.author.username}</span>
            <span>发布时间：${new Date(article.created_at).toLocaleDateString()}</span>
            <span>阅读：${article.view_count}</span>
            <span>评论：${article.comment_count}</span>
            <span>点赞：${article.like_count || 0}</span>
        </div>
        <div class="article-tags">
            ${article.tags && article.tags.length > 0 ? 
                article.tags.map(tag => `<span class="tag">${typeof tag === 'object' ? tag.name : tag}</span>`).join('') : 
                '无标签'
            }
        </div>
        <div class="article-actions">
            <button class="like-btn ${article.liked ? 'liked' : ''}" onclick="toggleLike(${article.id})">
                <span class="heart">${article.liked ? '❤️' : '🤍'}</span>
                <span class="like-count">${article.like_count || article.likes || 0}</span>
            </button>
        </div>
        <div class="article-body">
            ${article.content}
        </div>
    `;
}

// 加载评论
async function loadComments(articleId) {
    try {
        const response = await fetch(`${API_BASE_URL}/api/comments/article/${articleId}`);
        if (response.ok) {
            const comments = await response.json();
            displayComments(comments);
        }
    } catch (error) {
        console.error('加载评论失败：', error);
    }
}

// 显示评论
function displayComments(comments) {
    const container = document.getElementById('comments-list');
    container.innerHTML = '';

    if (comments.length === 0) {
        container.innerHTML = '<p>暂无评论</p>';
        return;
    }

    comments.forEach(comment => {
        const commentElement = document.createElement('div');
        commentElement.className = 'comment-item';
        
        commentElement.innerHTML = `
            <div class="comment-header">
                <span class="comment-author">${comment.user.username}</span>
                <span>${new Date(comment.created_at).toLocaleString()}</span>
            </div>
            <div class="comment-content">${comment.content}</div>
            <div class="comment-actions">
                <button onclick="likeComment(${comment.id})">点赞 (${comment.likes})</button>
            </div>
        `;
        
        container.appendChild(commentElement);
    });
}

// 提交评论
async function submitComment() {
    if (!currentUser || !currentArticleId) return;
    
    const content = document.getElementById('comment-content').value;
    if (!content.trim()) {
        alert('请输入评论内容');
        return;
    }

    try {
        const response = await fetch(`${API_BASE_URL}/api/comments`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${authToken}`,
            },
            body: JSON.stringify({
                content: content,
                article_id: currentArticleId,
            }),
        });

        if (response.ok) {
            document.getElementById('comment-content').value = '';
            loadComments(currentArticleId);
            alert('评论发表成功！');
        } else {
            alert('评论发表失败');
        }
    } catch (error) {
        alert('评论发表失败：网络错误');
    }
}

// 点赞评论
async function likeComment(commentId) {
    try {
        const response = await fetch(`${API_BASE_URL}/api/comments/${commentId}/like`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': authToken ? `Bearer ${authToken}` : ''
            }
        });

        if (response.ok) {
            loadComments(currentArticleId);
        }
    } catch (error) {
        console.error('点赞失败：', error);
    }
}

// 创建文章
async function createArticle() {
    if (!currentUser) {
        alert('请先登录');
        return;
    }

    // 直接获取表单数据 - 最简版本
    const title = document.getElementById('article-title').value.trim();
    const content = document.getElementById('article-content').value.trim();
    const summary = document.getElementById('article-summary').value.trim();
    const tagsText = document.getElementById('article-tags').value.trim();
    const status = document.getElementById('article-status').value;

    // 简单的验证
    if (!title) {
        alert('请填写文章标题');
        return;
    }
    if (!content) {
        alert('请填写文章内容');
        return;
    }

    const tags = tagsText ? tagsText.split(',').map(tag => tag.trim()).filter(tag => tag) : [];

    try {
        const requestData = {
            title,
            content,
            summary,
            tags,
            status,
        };
        
        const response = await fetch(`${API_BASE_URL}/api/articles`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${authToken}`,
            },
            body: JSON.stringify(requestData),
        });

        if (response.ok) {
            alert('文章保存成功！');
            showHome();
        } else {
            const errorData = await response.json();
            alert('文章保存失败：' + (errorData.error || '未知错误'));
        }
    } catch (error) {
        alert('文章保存失败：网络错误');
    }
}


// 加载个人资料
function loadProfile() {
    if (!currentUser) return;
    
    document.getElementById('profile-username').textContent = currentUser.username;
    document.getElementById('profile-email').textContent = currentUser.email;
    document.getElementById('profile-bio').textContent = currentUser.bio || '暂无个人简介';
    
    if (currentUser.avatar) {
        document.getElementById('profile-avatar').src = currentUser.avatar;
    }
}