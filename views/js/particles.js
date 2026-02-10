/**
 * Anime-style Particle Physics Engine
 * Features: Photon glow, constellation connections, mouse interaction, fluid physics
 */

(function () {
    const CONFIG = {
        particleCount: 120,
        connectionRadius: 120,
        mouseRadius: 180,
        colors: ['#ff79c6', '#7ab8ff', '#bd93f9', '#8be9fd', '#f8f8f2'], // Pink, Blue, Purple, Cyan, White
        baseSpeed: 0.8,
        glowBlur: 15,
        friction: 0.96, // Mouse push dampening
    };

    let canvas, ctx, w, h;
    let particles = [];
    let mouse = { x: null, y: null };

    class Particle {
        constructor() {
            this.x = Math.random() * w;
            this.y = Math.random() * h;
            this.vx = (Math.random() - 0.5) * CONFIG.baseSpeed;
            this.vy = (Math.random() - 0.5) * CONFIG.baseSpeed;
            this.size = Math.random() * 2 + 0.5;
            this.color = CONFIG.colors[Math.floor(Math.random() * CONFIG.colors.length)];
            this.alpha = Math.random() * 0.5 + 0.5;
            this.mx = 0; // Momentum from mouse
            this.my = 0;
        }

        update() {
            // Apply physics
            this.mx *= CONFIG.friction;
            this.my *= CONFIG.friction;

            this.x += this.vx + this.mx;
            this.y += this.vy + this.my;

            // Boundary wrap
            if (this.x < 0) this.x = w;
            if (this.x > w) this.x = 0;
            if (this.y < 0) this.y = h;
            if (this.y > h) this.y = 0;

            // Mouse interaction
            if (mouse.x != null) {
                let dx = mouse.x - this.x;
                let dy = mouse.y - this.y;
                let distance = Math.sqrt(dx * dx + dy * dy);

                if (distance < CONFIG.mouseRadius) {
                    const forceDirectionX = dx / distance;
                    const forceDirectionY = dy / distance;
                    const force = (CONFIG.mouseRadius - distance) / CONFIG.mouseRadius;
                    const directionMultiplier = 0.05; // Strength of push

                    this.mx -= forceDirectionX * force * directionMultiplier;
                    this.my -= forceDirectionY * force * directionMultiplier;
                }
            }
        }

        draw() {
            ctx.beginPath();
            ctx.arc(this.x, this.y, this.size, 0, Math.PI * 2);
            ctx.fillStyle = this.color;
            ctx.globalAlpha = this.alpha;
            ctx.shadowBlur = CONFIG.glowBlur; // Photon glow
            ctx.shadowColor = this.color;
            ctx.fill();
            ctx.shadowBlur = 0; // Reset for lines
        }
    }

    function init() {
        canvas = document.createElement("canvas");
        canvas.id = "anime-particles";
        canvas.style.position = "fixed";
        canvas.style.top = "0";
        canvas.style.left = "0";
        canvas.style.width = "100%";
        canvas.style.height = "100%";
        canvas.style.zIndex = "-1"; // Behind everything
        canvas.style.pointerEvents = "none";
        canvas.style.background = "radial-gradient(circle at center, #1b1b3a 0%, #0a0a1a 100%)"; // Deep space bg
        document.body.prepend(canvas);

        ctx = canvas.getContext("2d");

        resize();
        window.addEventListener("resize", resize);
        window.addEventListener("mousemove", (e) => {
            mouse.x = e.clientX;
            mouse.y = e.clientY;
        });
        window.addEventListener("mouseout", () => {
            mouse.x = null;
            mouse.y = null;
        });

        // Create particles
        for (let i = 0; i < CONFIG.particleCount; i++) {
            particles.push(new Particle());
        }

        animate();
    }

    function resize() {
        w = canvas.width = window.innerWidth;
        h = canvas.height = window.innerHeight;
    }

    function animate() {
        ctx.clearRect(0, 0, w, h);

        // Update and draw particles
        particles.forEach(p => {
            p.update();
            p.draw();
        });

        // Draw connections
        connect();

        requestAnimationFrame(animate);
    }

    function connect() {
        for (let a = 0; a < particles.length; a++) {
            for (let b = a; b < particles.length; b++) {
                let dx = particles[a].x - particles[b].x;
                let dy = particles[a].y - particles[b].y;
                let distance = Math.sqrt(dx * dx + dy * dy);

                if (distance < CONFIG.connectionRadius) {
                    let opacity = 1 - (distance / CONFIG.connectionRadius);
                    ctx.strokeStyle = `rgba(180, 200, 255, ${opacity * 0.4})`;
                    ctx.lineWidth = 1;
                    ctx.beginPath();
                    ctx.moveTo(particles[a].x, particles[a].y);
                    ctx.lineTo(particles[b].x, particles[b].y);
                    ctx.stroke();
                }
            }
        }
    }

    // Wait for DOM
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', init);
    } else {
        init();
    }
})();
