-- 清空现有数据
DELETE FROM comments;
DELETE FROM article_tags;
DELETE FROM articles;
DELETE FROM tags;
DELETE FROM categories;
DELETE FROM users;

-- 重置自增ID
DELETE FROM sqlite_sequence WHERE name IN ('users', 'categories', 'tags', 'articles', 'comments');

-- 插入用户数据
INSERT INTO users (username, email, password, avatar, bio, created_at, updated_at) VALUES 
('博主小明', 'admin@example.com', '$2a$10$jD1dYk5QO6aWxY3z5V3O0u5K5K5K5K5K5K5K5K5K5K5K5K5K5K5K5K5K', '/static/images/default-avatar.png', '热爱编程和分享的技术博主', datetime('now'), datetime('now')),
('技术达人', 'tech@example.com', '$2a$10$jD1dYk5QO6aWxY3z5V3O0u5K5K5K5K5K5K5K5K5K5K5K5K5K5K5K5K5K', '/static/images/default-avatar.png', '专注于后端开发的技术专家', datetime('now'), datetime('now')),
('前端小王', 'frontend@example.com', '$2a$10$jD1dYk5QO6aWxY3z5V3O0u5K5K5K5K5K5K5K5K5K5K5K5K5K5K5K5K5K', '/static/images/default-avatar.png', '前端开发工程师', datetime('now'), datetime('now'));

-- 插入分类数据
INSERT INTO categories (name, desc, user_id, created_at, updated_at) VALUES 
('生活随笔', '记录生活中的点点滴滴', 1, datetime('now'), datetime('now')),
('技术分享', '分享技术心得和学习笔记', 2, datetime('now'), datetime('now')),
('前端开发', '前端技术相关内容', 3, datetime('now'), datetime('now'));

-- 插入标签数据
INSERT INTO tags (name) VALUES 
('生活'),
('技术'),
('Go语言'),
('编程'),
('并发'),
('前端'),
('JavaScript'),
('最佳实践');

-- 插入文章数据
INSERT INTO articles (title, content, summary, author_id, category_id, status, view_count, comment_count, is_draft, created_at, updated_at) VALUES 
('欢迎来到我的博客', '这是我的第一篇博客文章！在这里我会分享我的技术心得、生活感悟和学习体会。希望能和大家一起成长，共同进步。作为一名技术爱好者，我一直在探索各种编程语言和开发工具。无论是前端还是后端，我都保持着浓厚的兴趣。通过这个博客平台，我希望能够记录我的学习历程，同时也为其他学习者提供一些有价值的参考资料。在未来的日子里，我会定期更新内容，包括技术教程、项目实践和行业动态等。期待与各位读者共同交流，共同成长！', '博客开篇，分享技术心得和生活感悟', 1, 1, 'published', 128, 5, false, datetime('now', '-10 day'), datetime('now', '-10 day')),
('Go语言学习笔记', '最近在学习Go语言，发现它的并发模型真的很强大。goroutine和channel的设计让并发编程变得简单易懂。Go语言由Google开发，于2009年正式发布，它的设计哲学是简单、高效、可靠。在现代软件开发中，Go语言越来越受到青睐，特别是在云原生、微服务、DevOps等领域。

在学习Go的过程中，我发现它有很多独特的特性。首先，Go的语法简洁明了，去除了很多传统语言中的冗余元素。其次，Go内置了垃圾回收机制，不需要手动管理内存。最重要的是，Go的并发模型非常强大，通过goroutine和channel，开发者可以轻松实现高效的并发程序。

Goroutine是Go语言中的轻量级线程，创建成本很低，一个程序可以创建成千上万个goroutine。而channel则是goroutine之间通信的桥梁，它确保了数据在不同goroutine之间安全传递。通过这两个核心概念，Go语言实现了"不要通过共享内存来通信，而要通过通信来共享内存"的设计理念。

除了并发特性外，Go语言还提供了丰富的标准库，涵盖了网络编程、文件操作、加密解密等各种功能。这使得开发者能够快速构建各种应用程序，而无需引入大量的第三方依赖。

总的来说，Go语言是一门非常适合现代软件开发的语言，特别是对于需要高并发、高性能的后端服务。如果你还没有尝试过Go语言，我强烈建议你抽时间学习一下，相信你会像我一样，被它的简洁和强大所吸引。', 'Go语言并发编程学习心得', 2, 2, 'published', 256, 8, false, datetime('now', '-5 day'), datetime('now', '-5 day')),
('前端开发最佳实践', '在前端开发中，有很多最佳实践值得我们遵循。比如组件化开发、状态管理、性能优化等。今天想和大家分享一些实用的技巧和经验。

首先，组件化开发是现代前端框架的核心思想。通过将UI拆分成独立的、可复用的组件，我们可以提高代码的可维护性和重用性。在实现组件时，应该遵循单一职责原则，确保每个组件只负责一个功能。同时，要注意组件的粒度，过大或过小的组件都会影响开发效率。

其次，状态管理是复杂前端应用的关键。随着应用规模的增长，组件间的状态共享和通信变得越来越复杂。这时，我们可以考虑使用Redux、Vuex或MobX等状态管理库。这些库提供了集中式的状态管理方案，使得状态的变化更加可预测和可追踪。

性能优化是前端开发中永恒的话题。常见的优化手段包括资源压缩、代码分割、懒加载、缓存策略等。在实际开发中，我们应该使用性能分析工具来识别瓶颈，然后有针对性地进行优化。记住，过早的优化是万恶之源，我们应该在确保功能正常的前提下进行优化。

另外，代码质量也是前端开发中不可忽视的一环。使用ESLint、Prettier等工具可以帮助我们保持一致的代码风格。编写单元测试和端到端测试可以提高代码的稳定性和可靠性。同时，定期进行代码审查也是提高团队代码质量的有效手段。

最后，持续学习是前端开发者必备的能力。前端技术更新迭代非常快，新框架、新工具层出不穷。我们应该保持开放的心态，不断学习新技术，同时也要深入理解前端开发的基础概念。只有这样，才能在快速变化的前端领域保持竞争力。

希望这些建议对大家有所帮助。如果有任何疑问或建议，欢迎在评论区留言讨论。', '前端开发的最佳实践和实用技巧', 3, 3, 'published', 189, 12, false, datetime('now', '-3 day'), datetime('now', '-3 day')),
('Docker容器化实践指南', 'Docker已经成为现代软件开发和部署的标准工具之一。容器化技术不仅简化了应用的部署过程，还提高了开发环境的一致性。在这篇文章中，我将分享一些Docker容器化的实践经验。

首先，让我们回顾一下Docker的基本概念。Docker使用容器来打包应用及其所有依赖，使其能够在任何环境中一致地运行。容器与虚拟机不同，它不需要模拟完整的操作系统，而是共享宿主机的内核，这使得容器更加轻量级和高效。

在实际项目中，编写良好的Dockerfile是容器化的第一步。一个好的Dockerfile应该遵循一些最佳实践，比如使用多阶段构建来减小镜像大小，使用适当的基础镜像，最小化镜像层数，以及设置非root用户运行等。

Next, let us talk about Docker Compose. For multi-container applications, Docker Compose provides a convenient way to define and run them. By writing a docker-compose.yml file, you can configure all the services your application needs, including databases, message queues, and web servers.

Another important aspect is image optimization. Large images not only consume more disk space but also slow down the deployment process. To optimize images, you can use techniques like layer caching, multi-stage builds, and Alpine-based images. Additionally, using a registry like Docker Hub or a private registry can help manage and distribute images efficiently.

在容器化过程中，安全性也是一个不可忽视的问题。我们应该定期更新基础镜像以修复已知漏洞，限制容器的权限，避免在容器中存储敏感信息，以及使用Docker的安全扫描工具来检测潜在的安全问题。

最后，对于生产环境的容器化部署，我们可以考虑使用Kubernetes等容器编排工具。Kubernetes提供了强大的容器编排能力，包括自动扩缩容、服务发现、负载均衡等功能，能够满足大规模应用的部署需求。

总之，Docker容器化是现代软件开发中的重要实践。通过合理使用Docker，我们可以提高开发效率，简化部署流程，确保环境一致性，从而更快地交付高质量的软件产品。', 'Docker容器化技术实践与经验分享', 2, 2, 'published', 320, 15, false, datetime('now', '-8 day'), datetime('now', '-8 day')),
('React Hooks深度解析', 'React Hooks是React 16.8引入的新特性，它允许你在不编写class的情况下使用state以及其他的React特性。Hooks彻底改变了React应用的开发方式，使得函数组件可以拥有类组件的所有功能。

Let me start by introducing the most commonly used hooks. useState allows functional components to manage state, useEffect handles side effects, useContext provides access to context, and useReducer offers an alternative to useState for complex state logic.

在使用React Hooks时，有一些规则需要遵守。首先，Hooks只能在函数组件的顶层调用，不能在循环、条件或嵌套函数中调用。其次，Hooks只能在React函数组件或自定义Hooks中调用。这些规则确保了Hooks在每次渲染时都能以相同的顺序被调用，从而保证了状态的正确保存。

自定义Hooks是Hooks的一个强大特性，它允许我们将组件逻辑提取到可重用的函数中。通过创建自定义Hooks，我们可以在不同的组件之间共享逻辑，避免重复代码，提高代码的可维护性。

在性能优化方面，React提供了一些专门的Hooks，如useMemo用于缓存计算结果，useCallback用于缓存函数引用。合理使用这些Hooks可以减少不必要的重渲染，提高应用的性能。

Another important aspect is testing components that use Hooks. React Testing Library provides utilities to test components with Hooks, allowing you to verify that your components behave correctly under different conditions.

随着React Hooks的广泛应用，社区中也出现了很多优秀的自定义Hooks库，如useHooks.com和react-use。这些库提供了大量现成的自定义Hooks，可以帮助我们快速实现各种功能，如网络请求、表单处理、状态管理等。

总之，React Hooks是React生态系统中的重要创新，它使得函数组件成为React开发的主流方式。通过掌握React Hooks，我们可以编写更加简洁、可维护的React代码，提高开发效率。

希望这篇文章能帮助你更好地理解和应用React Hooks。如果你有任何疑问或建议，欢迎在评论区留言讨论。', 'React Hooks特性详解与实战应用', 3, 3, 'published', 280, 20, false, datetime('now', '-6 day'), datetime('now', '-6 day')),
('数据分析与可视化入门指南', '在当今数字化时代，数据分析和可视化已经成为各行各业的重要技能。通过对数据的分析和可视化，我们可以发现数据中的模式和趋势，为决策提供支持。在这篇文章中，我将为大家介绍数据分析和可视化的基础知识和实用工具。

首先，让我们了解数据分析的基本流程。数据分析通常包括数据收集、数据清洗、数据探索、数据建模和结果解释等步骤。每个步骤都有其特定的方法和工具。

Data collection is the first step in the data analysis process. It involves gathering relevant data from various sources, such as databases, APIs, files, and sensors. The quality of the data collected directly affects the quality of the analysis results.

数据清洗是数据分析中最耗时但也最重要的步骤之一。原始数据往往存在缺失值、异常值和重复数据等问题，需要进行处理。常用的数据清洗工具包括Python的pandas库和R语言等。

数据探索是通过统计分析和可视化来了解数据特性的过程。在这个阶段，我们可以使用描述性统计指标（如均值、中位数、标准差等）来概括数据的基本特征，也可以使用各种图表（如直方图、散点图、箱线图等）来直观地展示数据的分布和关系。

数据建模是使用统计学或机器学习方法来分析数据并建立预测模型的过程。根据问题的性质，可以选择不同的建模方法，如线性回归、分类算法、聚类算法等。

数据可视化是将数据以图形化的方式呈现出来的过程。好的数据可视化能够直观地展示数据中的信息，帮助观众快速理解数据的含义。常用的数据可视化工具包括Tableau、Power BI、Python的matplotlib和seaborn库等。

在实际应用中，数据分析和可视化往往是结合在一起的。通过迭代地分析和可视化数据，我们可以不断深入了解数据，发现新的见解。

总之，数据分析和可视化是一项强大的技能，它可以帮助我们从数据中提取有价值的信息，支持决策制定。无论你是数据分析师、产品经理还是软件开发人员，掌握数据分析和可视化的基础知识都将对你的工作有所帮助。

希望这篇入门指南能够为你提供一些帮助。如果你对数据分析和可视化感兴趣，建议你进一步学习相关的理论知识，并通过实践来提高自己的技能。', '数据分析与可视化基础教程与工具推荐', 1, 1, 'published', 195, 8, false, datetime('now', '-4 day'), datetime('now', '-4 day')),
('微服务架构设计原则与实践', '微服务架构已经成为构建大型分布式系统的主流方法。与传统的单体架构相比，微服务架构将应用拆分为多个独立的服务，每个服务负责特定的业务功能。这种架构方式具有很多优势，如更好的可扩展性、更高的容错性、更灵活的技术栈选择等。

在设计微服务架构时，有一些关键原则需要遵循。首先是服务拆分原则，应该根据业务能力来划分服务，确保每个服务都是内聚的。其次是服务通信原则，微服务之间通常通过RESTful API、消息队列或gRPC等方式进行通信。此外，还有数据存储原则、服务发现与负载均衡原则、安全性原则等。

Service decomposition is a critical step in implementing a microservices architecture. There are several approaches to decomposition, including decomposition by business capability, decomposition by subdomain, and decomposition by team structure. The key is to ensure that each microservice has a clear responsibility and is loosely coupled with other services.

在服务通信方面，RESTful API是最常用的方式之一，它基于HTTP协议，简单易用。对于需要更高性能的场景，可以考虑使用gRPC，它基于Protocol Buffers，提供了更高效的序列化和通信机制。另外，对于异步通信场景，消息队列（如Kafka、RabbitMQ）也是一个很好的选择。

服务发现和负载均衡是微服务架构中的重要组件。在动态变化的微服务环境中，服务实例可能会频繁地启动和停止，因此需要一种机制来自动发现可用的服务实例并在它们之间分配流量。常用的服务发现工具包括Consul、Etcd和ZooKeeper等。

数据管理是微服务架构中的一个挑战。在微服务架构中，每个服务通常维护自己的数据存储，这带来了分布式事务、数据一致性等问题。为了解决这些问题，可以考虑使用Saga模式、CQRS模式等。

监控和可观测性对于微服务架构至关重要。由于微服务数量众多，系统复杂度高，因此需要强大的监控工具来跟踪服务的运行状态、性能指标和错误情况。常用的监控工具包括Prometheus、Grafana、ELK栈等。

在实际实施微服务架构时，还需要考虑很多因素，如团队结构、开发流程、部署策略等。微服务架构不仅是一种技术选择，更是一种组织和文化的转变。

总之，微服务架构为构建大型复杂系统提供了一种有效的方法，但它也带来了新的挑战。在采用微服务架构时，应该根据项目的实际情况，权衡利弊，选择合适的架构方案。', '微服务架构设计与实施指南', 2, 2, 'published', 240, 12, false, datetime('now', '-9 day'), datetime('now', '-9 day')),
('JavaScript异步编程全解析', 'JavaScript是一门单线程语言，但它通过异步编程模型实现了非阻塞操作。异步编程是JavaScript的核心特性之一，也是前端开发中必须掌握的技能。在这篇文章中，我将详细介绍JavaScript异步编程的各种方法和最佳实践。

首先，让我们回顾一下JavaScript异步编程的发展历程。从早期的回调函数（Callbacks），到Promise，再到async/await语法，JavaScript异步编程经历了几次重要的演进，变得越来越简洁和易用。

回调函数是JavaScript中最基本的异步编程方式。通过将一个函数作为参数传递给另一个函数，当异步操作完成时，调用该回调函数。然而，回调函数容易导致"回调地狱"（Callback Hell），代码嵌套层级过深，可读性和可维护性较差。

为了解决回调地狱问题，ES6引入了Promise。Promise是一种用于处理异步操作的对象，它代表了一个异步操作的最终完成（或失败）及其结果值。Promise提供了链式调用的方式，使得异步代码更加扁平，可读性更好。

ES2017引入了async/await语法，它基于Promise，提供了一种更接近同步代码的写法。通过使用async关键字定义异步函数，在函数内部使用await关键字等待Promise的解决，可以编写看起来像同步代码的异步代码。

Let us look at some practical examples. Suppose we need to fetch data from multiple APIs sequentially. Using callbacks, this would result in nested function calls. With Promises, we can use the then() method to chain the operations. And with async/await, we can write code that looks almost synchronous, making it much easier to read and understand.

在实际开发中，我们经常需要处理多个异步操作。对于并行执行的异步操作，可以使用Promise.all()、Promise.race()等方法。对于需要按顺序执行的异步操作，可以使用async/await配合for循环。

错误处理是异步编程中的重要环节。在Promise中，可以使用catch()方法捕获错误。在async/await中，可以使用try/catch语句捕获错误。合理的错误处理可以提高应用的健壮性。

在性能优化方面，需要注意避免不必要的异步操作嵌套，合理使用Promise.all()进行并行处理，以及注意内存泄漏问题。另外，对于长时间运行的任务，可以考虑使用Web Workers来避免阻塞主线程。

总之，JavaScript异步编程是前端开发中的重要技能。通过掌握各种异步编程方法，我们可以编写出更高效、更可靠的JavaScript代码。希望这篇文章能够帮助你更好地理解和应用JavaScript异步编程。', 'JavaScript异步编程方法与最佳实践', 3, 3, 'published', 310, 18, false, datetime('now', '-7 day'), datetime('now', '-7 day')),
('人工智能在日常生活中的应用', '人工智能（AI）技术正在快速发展，并逐渐渗透到我们日常生活的方方面面。从智能手机上的语音助手，到智能推荐系统，再到自动驾驶汽车，AI技术正在改变我们的生活方式。在这篇文章中，我将介绍人工智能在日常生活中的一些常见应用。

首先，让我们看一下AI在智能手机中的应用。现在的智能手机都配备了各种AI功能，如面部识别、语音助手、智能相机等。面部识别技术使用深度学习算法来识别人脸，用于解锁手机、支付验证等。语音助手如Siri、Google Assistant等，可以通过自然语言处理技术理解用户的指令并执行相应的操作。

智能推荐系统是AI的另一个广泛应用领域。在电商平台、视频网站、音乐流媒体服务等，推荐系统根据用户的历史行为和偏好，推荐个性化的商品、视频或音乐。这些推荐系统通常使用协同过滤、内容推荐等算法，通过分析大量用户数据来提供准确的推荐。

Another important application of AI is in healthcare. AI technologies are being used for medical imaging analysis, disease diagnosis, drug discovery, and personalized treatment. For example, AI algorithms can analyze medical images like X-rays and MRIs to detect abnormalities with high accuracy, helping doctors make better diagnoses.

智能家居是AI应用的另一个热门领域。智能音箱、智能灯光、智能温控器等设备，通过AI技术实现自动化控制和语音交互。这些设备可以根据用户的习惯自动调节家居环境，提高生活的舒适度和便利性。

在交通出行方面，AI技术也发挥着越来越重要的作用。自动驾驶汽车是AI技术的集大成者，它使用计算机视觉、机器学习等技术来感知周围环境，做出驾驶决策。虽然完全的自动驾驶技术还在发展中，但辅助驾驶功能如自动泊车、车道保持、自适应巡航等已经广泛应用于现代汽车中。

AI在教育领域也有很多应用。智能辅导系统可以根据学生的学习情况提供个性化的学习内容和辅导。AI批改系统可以自动批改作业和考试，减轻教师的工作负担。此外，AI还可以用于教育资源的推荐和学习路径的规划。

当然，AI技术的发展也带来了一些挑战，如隐私保护、就业影响、算法偏见等。这些问题需要我们在享受AI技术带来便利的同时，认真思考并寻找解决方案。

总之，人工智能技术正在深刻改变我们的生活方式。随着AI技术的不断发展，相信它将在更多领域发挥重要作用，为我们的生活带来更多便利和创新。', '人工智能技术在日常生活中的应用与影响', 1, 1, 'published', 265, 14, false, datetime('now', '-2 day'), datetime('now', '-2 day')),

('云计算平台比较：AWS、Azure与GCP', '云计算已经成为现代IT基础设施的重要组成部分。三大云服务提供商——Amazon Web Services（AWS）、Microsoft Azure和Google Cloud Platform（GCP）占据了云计算市场的主导地位。在这篇文章中，我将对这三大云平台进行比较，帮助你了解它们的特点和优势。

首先，让我们看一下AWS。作为最早进入云计算市场的玩家之一，AWS拥有最成熟和全面的服务生态系统。它提供了丰富的计算、存储、数据库、网络、分析、人工智能等服务，可以满足各种规模企业的需求。AWS的全球基础设施覆盖广泛，在全球多个区域和可用区部署了数据中心，提供了高可用性和低延迟的服务。

Microsoft Azure是微软推出的云平台，它与微软的企业软件产品线（如Windows Server、SQL Server、Office 365等）有着深度集成。对于已经使用微软产品的企业来说，迁移到Azure可以获得更好的兼容性和集成体验。Azure在混合云领域有着独特优势，提供了丰富的混合连接选项和工具。

Google Cloud Platform是Google推出的云平台，它在大数据分析、机器学习和容器化方面有着强大的优势。GCP的BigQuery是一个强大的无服务器数据仓库，Cloud ML Engine提供了简单易用的机器学习服务，而Kubernetes（Google开源并维护）在GCP上有着原生支持。此外，GCP的网络基础设施基于Google的全球网络，提供了优质的网络性能。

Compare these three platforms in terms of pricing. Pricing models vary across services and providers, but generally speaking, all three offer pay-as-you-go pricing with various discounts for reserved instances or commitments. AWS tends to have the most complex pricing structure with many options and discounts. Azure often provides price matching guarantees against AWS for comparable services. GCP typically offers sustained use discounts that automatically apply based on usage patterns.

在安全性方面，三大云平台都提供了多层次的安全防护措施，包括网络安全、身份认证、数据加密、合规认证等。AWS的Security Hub和GuardDuty提供了全面的安全监控和威胁检测功能。Azure的Security Center提供了统一的安全管理视图。GCP的Security Command Center整合了各种安全服务，提供了集中式的安全管理。

在开发者体验方面，三大云平台都提供了完善的开发工具和SDK。AWS的AWS CLI和SDK支持多种编程语言，提供了丰富的API。Azure的Azure CLI和SDK与Visual Studio和Visual Studio Code有着良好的集成。GCP的Cloud SDK和Cloud Code插件为开发者提供了便捷的开发环境。

选择云平台时，需要考虑多种因素，如业务需求、现有技术栈、预算、合规要求等。对于初创企业或需要快速扩展的业务，AWS的成熟生态和灵活服务可能是不错的选择。对于已经深度使用微软产品的企业，Azure的无缝集成可能更有吸引力。对于需要强大数据处理和机器学习能力的业务，GCP可能更合适。

总之，AWS、Azure和GCP各有优势，选择哪个云平台应该基于企业的具体需求和场景。在某些情况下，混合使用多个云平台，充分利用各平台的优势，可能是更好的策略。

希望这篇比较文章能够帮助你更好地了解三大云平台。如果你对云计算感兴趣，建议你进一步学习各平台的具体服务和最佳实践。', '三大主流云服务平台对比分析', 2, 2, 'published', 290, 16, false, datetime('now', '-1 day'), datetime('now', '-1 day'));

-- 插入文章标签关联数据
INSERT INTO article_tags (article_id, tag_id) VALUES 
(1, 1),
(1, 2),
(2, 3),
(2, 4),
(2, 5),
(3, 6),
(3, 7),
(3, 8),
(4, 2),
(4, 4),
(4, 8),
(5, 6),
(5, 7),
(5, 8),
(6, 1),
(6, 2),
(6, 4),
(7, 2),
(7, 4),
(7, 8),
(8, 6),
(8, 7),
(8, 4),
(9, 1),
(9, 2),
(9, 4),
(10, 2),
(10, 4),
(10, 8);

-- 插入评论数据
INSERT INTO comments (content, article_id, user_id, likes, created_at, updated_at) VALUES 
('写得很好，期待更多分享！', 1, 2, 5, datetime('now', '-10 day', '+1 hour'), datetime('now', '-10 day', '+1 hour')),
('博客界面很简洁，喜欢这种风格。', 1, 3, 3, datetime('now', '-10 day', '+2 hour'), datetime('now', '-10 day', '+2 hour')),
('Go语言确实很适合并发编程！', 2, 1, 8, datetime('now', '-5 day', '+1 hour'), datetime('now', '-5 day', '+1 hour')),
('这些技巧很实用，已经用上了！', 3, 1, 12, datetime('now', '-3 day', '+1 hour'), datetime('now', '-3 day', '+1 hour')),
('Docker确实让部署变得简单了很多！', 4, 1, 6, datetime('now', '-8 day', '+2 hour'), datetime('now', '-8 day', '+2 hour')),
('容器化是现代开发的必备技能', 4, 3, 4, datetime('now', '-8 day', '+3 hour'), datetime('now', '-8 day', '+3 hour')),
('Hooks让React代码更简洁了', 5, 2, 9, datetime('now', '-6 day', '+1 hour'), datetime('now', '-6 day', '+1 hour')),
('这篇文章对React Hooks解释得很清楚', 5, 1, 7, datetime('now', '-6 day', '+2 hour'), datetime('now', '-6 day', '+2 hour')),
('数据分析确实很重要', 6, 2, 5, datetime('now', '-4 day', '+3 hour'), datetime('now', '-4 day', '+3 hour')),
('谢谢分享这些可视化工具', 6, 3, 3, datetime('now', '-4 day', '+4 hour'), datetime('now', '-4 day', '+4 hour')),
('微服务架构需要谨慎实施', 7, 1, 8, datetime('now', '-9 day', '+2 hour'), datetime('now', '-9 day', '+2 hour')),
('服务拆分是关键', 7, 3, 6, datetime('now', '-9 day', '+3 hour'), datetime('now', '-9 day', '+3 hour')),
('异步编程确实是JavaScript的核心', 8, 2, 10, datetime('now', '-7 day', '+1 hour'), datetime('now', '-7 day', '+1 hour')),
('async/await让异步代码可读性好多了', 8, 1, 7, datetime('now', '-7 day', '+2 hour'), datetime('now', '-7 day', '+2 hour')),
('AI技术发展真快', 9, 3, 8, datetime('now', '-2 day', '+1 hour'), datetime('now', '-2 day', '+1 hour')),
('期待AI在更多领域的应用', 9, 2, 6, datetime('now', '-2 day', '+2 hour'), datetime('now', '-2 day', '+2 hour')),
('云平台选择需要根据实际需求', 10, 1, 9, datetime('now', '-1 day', '+3 hour'), datetime('now', '-1 day', '+3 hour')),
('三大云平台各有优势', 10, 3, 5, datetime('now', '-1 day', '+4 hour'), datetime('now', '-1 day', '+4 hour'));

-- 设置ID自增
UPDATE sqlite_sequence SET seq = 3 WHERE name = 'users';
UPDATE sqlite_sequence SET seq = 3 WHERE name = 'categories';
UPDATE sqlite_sequence SET seq = 8 WHERE name = 'tags';
UPDATE sqlite_sequence SET seq = 10 WHERE name = 'articles';
UPDATE sqlite_sequence SET seq = 18 WHERE name = 'comments';
