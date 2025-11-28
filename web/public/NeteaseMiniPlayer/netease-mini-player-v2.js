/**
 * [NMPv2] NeteaseMiniPlayer v2 JavaScript
 * Lightweight Player Component Based on NetEase Cloud Music API
 * 
 * Copyright 2025 BHCN STUDIO & 北海的佰川（ImBHCN[numakkiyu]）
 * 
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 * 
 *     http://www.apache.org/licenses/LICENSE-2.0
 * 
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
(()=>{try{const s=document.currentScript;if(s&&s.src){fetch(s.src,{mode:'cors',credentials:'omit'}).catch(()=>{});}}catch(e){}})();
const GlobalAudioManager = {
    currentPlayer: null,
    setCurrent(player) {
        if (this.currentPlayer && this.currentPlayer !== player) {
            this.currentPlayer.pause();
        }
        this.currentPlayer = player;
    }
};

const ICONS = {
    prev: `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 640 640"><path d="M556.2 541.6C544.2 546.6 530.5 543.8 521.3 534.7L352 365.3L352 512C352 524.9 344.2 536.6 332.2 541.6C320.2 546.6 306.5 543.8 297.3 534.7L128 365.3L128 512C128 529.7 113.7 544 96 544C78.3 544 64 529.7 64 512L64 128C64 110.3 78.3 96 96 96C113.7 96 128 110.3 128 128L128 274.7L297.4 105.4C306.6 96.2 320.3 93.5 332.3 98.5C344.3 103.5 352 115.1 352 128L352 274.7L521.4 105.3C530.6 96.1 544.3 93.4 556.3 98.4C568.3 103.4 576 115.1 576 128L576 512C576 524.9 568.2 536.6 556.2 541.6z"/></svg>`,
    next: `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 640 640"><path d="M83.8 541.6C95.8 546.6 109.5 543.8 118.7 534.7L288 365.3L288 512C288 524.9 295.8 536.6 307.8 541.6C319.8 546.6 333.5 543.8 342.7 534.7L512 365.3L512 512C512 529.7 526.3 544 544 544C561.7 544 576 529.7 576 512L576 128C576 110.3 561.7 96 544 96C526.3 96 512 110.3 512 128L512 274.7L342.6 105.3C333.4 96.1 319.7 93.4 307.7 98.4C295.7 103.4 288 115.1 288 128L288 274.7L118.6 105.4C109.4 96.2 95.7 93.5 83.7 98.5C71.7 103.5 64 115.1 64 128L64 512C64 524.9 71.8 536.6 83.8 541.6z"/></svg>`,
    play: `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 640 640"><path d="M187.2 100.9C174.8 94.1 159.8 94.4 147.6 101.6C135.4 108.8 128 121.9 128 136L128 504C128 518.1 135.5 531.2 147.6 538.4C159.7 545.6 174.8 545.9 187.2 539.1L523.2 355.1C536 348.1 544 334.6 544 320C544 305.4 536 291.9 523.2 284.9L187.2 100.9z"/></svg>`,
    pause: `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 640 640"><path d="M176 96C149.5 96 128 117.5 128 144L128 496C128 522.5 149.5 544 176 544L240 544C266.5 544 288 522.5 288 496L288 144C288 117.5 266.5 96 240 96L176 96zM400 96C373.5 96 352 117.5 352 144L352 496C352 522.5 373.5 544 400 544L464 544C490.5 544 512 522.5 512 496L512 144C512 117.5 490.5 96 464 96L400 96z"/></svg>`,
    volume: `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 640 640"><path d="M533.6 96.5C523.3 88.1 508.2 89.7 499.8 100C491.4 110.3 493 125.4 503.3 133.8C557.5 177.8 592 244.8 592 320C592 395.2 557.5 462.2 503.3 506.3C493 514.7 491.5 529.8 499.8 540.1C508.1 550.4 523.3 551.9 533.6 543.6C598.5 490.7 640 410.2 640 320C640 229.8 598.5 149.2 533.6 96.5zM473.1 171C462.8 162.6 447.7 164.2 439.3 174.5C430.9 184.8 432.5 199.9 442.8 208.3C475.3 234.7 496 274.9 496 320C496 365.1 475.3 405.3 442.8 431.8C432.5 440.2 431 455.3 439.3 465.6C447.6 475.9 462.8 477.4 473.1 469.1C516.3 433.9 544 380.2 544 320.1C544 260 516.3 206.3 473.1 171.1zM412.6 245.5C402.3 237.1 387.2 238.7 378.8 249C370.4 259.3 372 274.4 382.3 282.8C393.1 291.6 400 305 400 320C400 335 393.1 348.4 382.3 357.3C372 365.7 370.5 380.8 378.8 391.1C387.1 401.4 402.3 402.9 412.6 394.6C434.1 376.9 448 350.1 448 320C448 289.9 434.1 263.1 412.6 245.5zM80 416L128 416L262.1 535.2C268.5 540.9 276.7 544 285.2 544C304.4 544 320 528.4 320 509.2L320 130.8C320 111.6 304.4 96 285.2 96C276.7 96 268.5 99.1 262.1 104.8L128 224L80 224C53.5 224 32 245.5 32 272L32 368C32 394.5 53.5 416 80 416z"/></svg>`,
    lyrics: `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 640 640"><path d="M532 71C539.6 77.1 544 86.3 544 96L544 400C544 444.2 501 480 448 480C395 480 352 444.2 352 400C352 355.8 395 320 448 320C459.2 320 470 321.6 480 324.6L480 207.9L256 257.7L256 464C256 508.2 213 544 160 544C107 544 64 508.2 64 464C64 419.8 107 384 160 384C171.2 384 182 385.6 192 388.6L192 160C192 145 202.4 132 217.1 128.8L505.1 64.8C514.6 62.7 524.5 65 532.1 71.1z"/></svg>`,
    list: `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 640 640"><path d="M104 112C90.7 112 80 122.7 80 136L80 184C80 197.3 90.7 208 104 208L152 208C165.3 208 176 197.3 176 184L176 136C176 122.7 165.3 112 152 112L104 112zM256 128C238.3 128 224 142.3 224 160C224 177.7 238.3 192 256 192L544 192C561.7 192 576 177.7 576 160C576 142.3 561.7 128 544 128L256 128zM256 288C238.3 288 224 302.3 224 320C224 337.7 238.3 352 256 352L544 352C561.7 352 576 337.7 576 320C576 302.3 561.7 288 544 288L256 288zM256 448C238.3 448 224 462.3 224 480C224 497.7 238.3 512 256 512L544 512C561.7 512 576 497.7 576 480C576 462.3 561.7 448 544 448L256 448zM80 296L80 344C80 357.3 90.7 368 104 368L152 368C165.3 368 176 357.3 176 344L176 296C176 282.7 165.3 272 152 272L104 272C90.7 272 80 282.7 80 296zM104 432C90.7 432 80 442.7 80 456L80 504C80 517.3 90.7 528 104 528L152 528C165.3 528 176 517.3 176 504L176 456C176 442.7 165.3 432 152 432L104 432z"/></svg>`,
    minimize: `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 640 640"><path d="M64 320C64 178.6 178.6 64 320 64C461.4 64 576 178.6 576 320C576 461.4 461.4 576 320 576C178.6 576 64 461.4 64 320zM320 352C302.3 352 288 337.7 288 320C288 302.3 302.3 288 320 288C337.7 288 352 302.3 352 320C352 337.7 337.7 352 320 352zM224 320C224 373 267 416 320 416C373 416 416 373 416 320C416 267 373 224 320 224C267 224 224 267 224 320zM168 304C168 271.6 184.3 237.4 210.8 210.8C237.3 184.2 271.6 168 304 168C317.3 168 328 157.3 328 144C328 130.7 317.3 120 304 120C256.1 120 210.3 143.5 176.9 176.9C143.5 210.3 120 256.1 120 304C120 317.3 130.7 328 144 328C157.3 328 168 317.3 168 304z"/></svg>`, 
    maximize: `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 640 640"><path d="M64 320C64 178.6 178.6 64 320 64C461.4 64 576 178.6 576 320C576 461.4 461.4 576 320 576C178.6 576 64 461.4 64 320zM320 352C302.3 352 288 337.7 288 320C288 302.3 302.3 288 320 288C337.7 288 352 302.3 352 320C352 337.7 337.7 352 320 352zM224 320C224 373 267 416 320 416C373 416 416 373 416 320C416 267 373 224 320 224C267 224 224 267 224 320zM168 304C168 271.6 184.3 237.4 210.8 210.8C237.3 184.2 271.6 168 304 168C317.3 168 328 157.3 328 144C328 130.7 317.3 120 304 120C256.1 120 210.3 143.5 176.9 176.9C143.5 210.3 120 256.1 120 304C120 317.3 130.7 328 144 328C157.3 328 168 317.3 168 304z"/></svg>`, 
    loopList: `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 640 640"><path d="M534.6 182.6C547.1 170.1 547.1 149.8 534.6 137.3L470.6 73.3C461.4 64.1 447.7 61.4 435.7 66.4C423.7 71.4 416 83.1 416 96L416 128L256 128C150 128 64 214 64 320C64 337.7 78.3 352 96 352C113.7 352 128 337.7 128 320C128 249.3 185.3 192 256 192L416 192L416 224C416 236.9 423.8 248.6 435.8 253.6C447.8 258.6 461.5 255.8 470.7 246.7L534.7 182.7zM105.4 457.4C92.9 469.9 92.9 490.2 105.4 502.7L169.4 566.7C178.6 575.9 192.3 578.6 204.3 573.6C216.3 568.6 224 556.9 224 544L224 512L384 512C490 512 576 426 576 320C576 302.3 561.7 288 544 288C526.3 288 512 302.3 512 320C512 390.7 454.7 448 384 448L224 448L224 416C224 403.1 216.2 391.4 204.2 386.4C192.2 381.4 178.5 384.2 169.3 393.3L105.3 457.3z"/></svg>`,
    loopSingle: `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 640 640"><path d="M534.6 182.6C547.1 170.1 547.1 149.8 534.6 137.3L470.6 73.3C461.4 64.1 447.7 61.4 435.7 66.4C423.7 71.4 416 83.1 416 96L416 128L256 128C150 128 64 214 64 320C64 337.7 78.3 352 96 352C113.7 352 128 337.7 128 320C128 249.3 185.3 192 256 192L416 192L416 224C416 236.9 423.8 248.6 435.8 253.6C447.8 258.6 461.5 255.8 470.7 246.7L534.7 182.7zM105.4 457.4C92.9 469.9 92.9 490.2 105.4 502.7L169.4 566.7C178.6 575.9 192.3 578.6 204.3 573.6C216.3 568.6 224 556.9 224 544L224 512L384 512C490 512 576 426 576 320C576 302.3 561.7 288 544 288C526.3 288 512 302.3 512 320C512 390.7 454.7 448 384 448L224 448L224 416C224 403.1 216.2 391.4 204.2 386.4C192.2 381.4 178.5 384.2 169.3 393.3L105.3 457.3z"/><path d="M295 280L305 260L335 260L335 380L305 380L305 280Z"/></svg>`, 
    shuffle: `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 640 640"><path d="M467.8 98.4C479.8 93.4 493.5 96.2 502.7 105.3L566.7 169.3C572.7 175.3 576.1 183.4 576.1 191.9C576.1 200.4 572.7 208.5 566.7 214.5L502.7 278.5C493.5 287.7 479.8 290.4 467.8 285.4C455.8 280.4 448 268.9 448 256L448 224L416 224C405.9 224 396.4 228.7 390.4 236.8L358 280L318 226.7L339.2 198.4C357.3 174.2 385.8 160 416 160L448 160L448 128C448 115.1 455.8 103.4 467.8 98.4zM218 360L258 413.3L236.8 441.6C218.7 465.8 190.2 480 160 480L96 480C78.3 480 64 465.7 64 448C64 430.3 78.3 416 96 416L160 416C170.1 416 179.6 411.3 185.6 403.2L218 360zM502.6 534.6C493.4 543.8 479.7 546.5 467.7 541.5C455.7 536.5 448 524.9 448 512L448 480L416 480C385.8 480 357.3 465.8 339.2 441.6L185.6 236.8C179.6 228.7 170.1 224 160 224L96 224C78.3 224 64 209.7 64 192C64 174.3 78.3 160 96 160L160 160C190.2 160 218.7 174.2 236.8 198.4L390.4 403.2C396.4 411.3 405.9 416 416 416L448 416L448 384C448 371.1 455.8 359.4 467.8 354.4C479.8 349.4 493.5 352.2 502.7 361.3L566.7 425.3C572.7 431.3 576.1 439.4 576.1 447.9C576.1 456.4 572.7 464.5 566.7 470.5L502.7 534.5z"/></svg>`
};
class NeteaseMiniPlayer {
    constructor(element) {
        this.element = element;
        this.element.neteasePlayer = this;
        this.config = this.parseConfig();
        this.currentSong = null;
        this.playlist = [];
        this.currentIndex = 0;
        this.audio = new Audio();
        this.wasPlayingBeforeHidden = false;
        this.isPlaying = false;
        this.currentTime = 0;
        this.duration = 0;
        this.volume = 0.7;
        this.lyrics = [];
        this.currentLyricIndex = -1;
        this.showLyrics = this.config.lyric;
        this.cache = new Map();
        this.init();
        this.playMode = 'list';
        this.shuffleHistory = [];
        this.idleTimeout = null;
        this.idleDelay = 5000;
        this.isIdle = false;
    }
    parseConfig() {
        const element = this.element;
        const position = element.dataset.position || 'static';
        const validPositions = ['static', 'top-left', 'top-right', 'bottom-left', 'bottom-right'];
        const finalPosition = validPositions.includes(position) ? position : 'static';
        const defaultMinimized = element.dataset.defaultMinimized === 'true';
        
        const embedValue = element.getAttribute('data-embed') || element.dataset.embed;
        const isEmbed = embedValue === 'true' || embedValue === true;

        const autoPauseAttr = element.getAttribute('data-auto-pause') ?? element.dataset.autoPause;
        const autoPauseDisabled = autoPauseAttr === 'true' || autoPauseAttr === true;

        return {
            embed: isEmbed,
            autoplay: element.dataset.autoplay === 'true',
            playlistId: element.dataset.playlistId,
            songId: element.dataset.songId,
            position: finalPosition,
            lyric: element.dataset.lyric !== 'false',
            theme: element.dataset.theme || 'auto',
            size: element.dataset.size || 'compact',
            defaultMinimized: defaultMinimized,
            autoPauseDisabled: autoPauseDisabled
        };
    }
    async init() {
        if (this.config.embed) {
            this.element.setAttribute('data-embed', 'true');
        }
        this.element.setAttribute('data-position', this.config.position);
        
        if (this.config.embed) {
            this.element.classList.add('netease-mini-player-embed');
        }
        
        this.initTheme();
        this.createPlayerHTML();
        this.applyResponsiveControls?.();
        this.setupEnvListeners?.();
        this.bindEvents();
        this.setupAudioEvents();
        try {
            if (this.config.embed) {
                if (this.config.songId) {
                    await this.loadSingleSong(this.config.songId);
                } else if (this.config.playlistId) {
                    await this.loadPlaylist(this.config.playlistId);
                    this.playlist = [this.playlist[0]];
                }
            } else {
                if (this.config.playlistId) {
                    await this.loadPlaylist(this.config.playlistId);
                } else if (this.config.songId) {
                    await this.loadSingleSong(this.config.songId);
                }
            }
            if (this.playlist.length > 0) {
                await this.loadCurrentSong();
                if (this.config.autoplay && !this.config.embed) {
                    this.play();
                }
            }
            if (this.config.defaultMinimized && !this.config.embed && this.config.position !== 'static') {
                this.toggleMinimize();
            }
        } catch (error) {
            console.error('播放器初始化失败:', error);
            this.showError('加载失败，请稍后重试');
        }
    }
    createPlayerHTML() {
        this.element.innerHTML = `
            <div class="player-main">
                <div class="album-cover-container">
                    <img class="album-cover" src="" alt="专辑封面">
                    <div class="vinyl-overlay">
                        <div class="vinyl-center"></div>
                    </div>
                </div>
                <div class="song-content">
                    <div class="song-info">
                        <div class="song-title">加载中...</div>
                        <div class="song-artist">请稍候</div>
                    </div>
                    <div class="lyrics-container">
                        <div class="lyric-line original">♪ 加载歌词中... ♪</div>
                        <div class="lyric-line translation"></div>
                    </div>
                </div>
                <div class="controls">
                    ${!this.config.embed ? `<button class="control-btn prev-btn" title="上一首">${ICONS.prev}</button>` : ''}
                    <button class="control-btn play-btn" title="播放/暂停">
                        <span class="play-icon">${ICONS.play}</span>
                        <span class="pause-icon" style="display: none;">${ICONS.pause}</span>
                    </button>
                    ${!this.config.embed ? `<button class="control-btn next-btn" title="下一首">${ICONS.next}</button>` : ''}
                </div>
            </div>
            <div class="player-bottom">
                <div class="progress-container">
                    <span class="time-display current-time">0:00</span>
                    <div class="progress-bar-container">
                        <div class="progress-bar"></div>
                    </div>
                    <span class="time-display total-time">0:00</span>
                </div>
                <div class="bottom-controls">
                    <div class="volume-container">
                        <span class="volume-icon">${ICONS.volume}</span>
                        <div class="volume-slider-container">
                            <div class="volume-slider">
                                <div class="volume-bar"></div>
                            </div>
                        </div>
                    </div>
                    <span class="feature-btn lyrics-btn" role="button" title="显示/隐藏歌词">${ICONS.lyrics}</span>
                    ${!this.config.embed ? `<span class="feature-btn loop-mode-btn" role="button" title="列表循环">${ICONS.loopList}</span>` : ''}
                    ${!this.config.embed ? `<span class="feature-btn list-btn" role="button" title="播放列表">${ICONS.list}</span>` : ''}
                    ${!this.config.embed ? `<span class="feature-btn minimize-btn" role="button" title="缩小/展开">${ICONS.minimize}</span>` : ''}
                </div>
            </div>
            <div class="playlist-container">
                <div class="playlist-content"></div>
            </div>
        `;
        this.elements = {
            albumCover: this.element.querySelector('.album-cover'),
            albumCoverContainer: this.element.querySelector('.album-cover-container'),
            songTitle: this.element.querySelector('.song-title'),
            songArtist: this.element.querySelector('.song-artist'),
            lyricsContainer: this.element.querySelector('.lyrics-container'),
            lyricLine: this.element.querySelector('.lyric-line.original'),
            lyricTranslation: this.element.querySelector('.lyric-line.translation'),
            playBtn: this.element.querySelector('.play-btn'),
            playIcon: this.element.querySelector('.play-icon'),
            pauseIcon: this.element.querySelector('.pause-icon'),
            prevBtn: this.element.querySelector('.prev-btn'),
            nextBtn: this.element.querySelector('.next-btn'),
            progressContainer: this.element.querySelector('.progress-bar-container'),
            progressBar: this.element.querySelector('.progress-bar'),
            currentTime: this.element.querySelector('.current-time'),
            totalTime: this.element.querySelector('.total-time'),
            volumeContainer: this.element.querySelector('.volume-container'),
            volumeSlider: this.element.querySelector('.volume-slider'),
            volumeBar: this.element.querySelector('.volume-bar'),
            volumeIcon: this.element.querySelector('.volume-icon'),
            lyricsBtn: this.element.querySelector('.lyrics-btn'),
            listBtn: this.element.querySelector('.list-btn'),
            minimizeBtn: this.element.querySelector('.minimize-btn'),
            playlistContainer: this.element.querySelector('.playlist-container'),
            playlistContent: this.element.querySelector('.playlist-content')
        };
        this.isMinimized = false;
        this.elements.loopModeBtn = this.element.querySelector('.loop-mode-btn');
    }
    bindEvents() {
        this.elements.playBtn.addEventListener('click', () => this.togglePlay());
        if (this.elements.prevBtn) {
            this.elements.prevBtn.addEventListener('click', () => this.previousSong());
        }
        if (this.elements.nextBtn) {
            this.elements.nextBtn.addEventListener('click', () => this.nextSong());
        }
        if (this.elements.loopModeBtn) {
            this.elements.loopModeBtn.addEventListener('click', () => this.togglePlayMode());
        }
        this.elements.albumCoverContainer.addEventListener('click', () => {
            if (this.element.classList.contains('minimized')) {
                this.elements.albumCoverContainer.classList.toggle('expanded');
                return;
            }
            if (this.currentSong && this.currentSong.id) {
            const songUrl = `https://music.163.com/song?id=${this.currentSong.id}`;
            window.open(songUrl, '_blank', 'noopener,noreferrer');
            }
        });
        let isDragging = false;
        this.elements.progressContainer.addEventListener('mousedown', (e) => {
            isDragging = true;
            this.seekTo(e);
        });
        document.addEventListener('mousemove', (e) => {
            if (isDragging) {
                this.seekTo(e);
            }
        });
        document.addEventListener('mouseup', () => {
            isDragging = false;
        });
        this.elements.progressContainer.addEventListener('click', (e) => this.seekTo(e));
        let isVolumesDragging = false;
        this.elements.volumeSlider.addEventListener('mousedown', (e) => {
            isVolumesDragging = true;
            this.setVolume(e);
        });
        document.addEventListener('mousemove', (e) => {
            if (isVolumesDragging) {
                this.setVolume(e);
            }
        });
        document.addEventListener('mouseup', () => {
            isVolumesDragging = false;
        });
        this.elements.volumeSlider.addEventListener('click', (e) => this.setVolume(e));
        this.elements.volumeIcon.addEventListener('click', () => {
            if (this.audio.volume > 0) {
                this.audio.volume = 0;
                this.volume = 0;
            } else {
                this.audio.volume = 0.7;
                this.volume = 0.7;
            }
            this.updateVolumeDisplay();
        });
        this.elements.lyricsBtn.addEventListener('click', () => this.toggleLyrics());
        if (this.elements.listBtn) {
            this.elements.listBtn.addEventListener('click', () => this.togglePlaylist());
        }
        if (this.elements.minimizeBtn) {
            this.elements.minimizeBtn.addEventListener('click', () => this.toggleMinimize());
        }
        this.element.addEventListener('mouseenter', () => {
            if (this.isIdle && this.isMinimized) {
                this.popOut();
            }
        });
        this.element.addEventListener('mouseleave', () => {
            if (this.isMinimized) {
                this.dock();
            }
        });
    }
    setupAudioEvents() {
        this.audio.addEventListener('timeupdate', () => {
            this.currentTime = this.audio.currentTime;
            this.updateProgressDisplay();
            this.updateCurrentLyric();
        });
        this.audio.addEventListener('durationchange', () => {
            this.duration = this.audio.duration;
            this.updateProgressDisplay();
        });
        this.audio.addEventListener('play', () => {
            this.isPlaying = true;
            GlobalAudioManager.setCurrent(this);
            this.element.classList.add('player-playing');
            this.elements.albumCover.classList.add('playing');
        });
        this.audio.addEventListener('pause', () => {
            this.isPlaying = false;
            if (GlobalAudioManager.currentPlayer === this) {
                GlobalAudioManager.currentPlayer = null;
            }
            this.element.classList.remove('player-playing');
            this.elements.albumCover.classList.remove('playing');
        });
        this.audio.addEventListener('ended', () => this.handleSongEnd());
        document.addEventListener('visibilitychange', () => {
            if (document.hidden && !this.config.autoPauseDisabled) {
                this.wasPlayingBeforeHidden = !this.audio.paused;
                this.audio.pause();
            } else if (!document.hidden && this.wasPlayingBeforeHidden) {
                this.audio.play().catch(() => {});
                this.wasPlayingBeforeHidden = false;
            }
        });
    }
    async loadPlaylist(playlistId) {
        try {
            this.elements.songTitle.textContent = '加载歌单...';
            this.elements.songArtist.textContent = '';
            const response = await fetch(`https://api.hypcvgm.top/nmp/v2/playlist?id=${encodeURIComponent(playlistId)}`);
            if (!response.ok) throw new Error('获取歌单失败');
            const data = await response.json();
            if (!Array.isArray(data.tracks)) throw new Error('歌单数据格式错误');
            this.playlist = data.tracks.map(track => ({
                id: track.id,
                name: track.name,
                artist: track.ar?.map(a => a.name).join(', ') || track.artist || '未知歌手',
                cover: track.al?.picUrl || track.cover || '',
                url: track.url || ''
            }));
            this.updatePlaylistDisplay();
        } catch (error) {
            console.error('加载歌单失败:', error);
            this.showError('加载歌单失败');
        }
    }
    async loadSingleSong(songId) {
        try {
            this.elements.songTitle.textContent = '加载单曲...';
            this.elements.songArtist.textContent = '';
            const response = await fetch(`https://api.hypcvgm.top/nmp/v2/song?id=${encodeURIComponent(songId)}`);
            if (!response.ok) throw new Error('获取歌曲失败');
            const data = await response.json();
            const song = {
                id: data.id,
                name: data.name,
                artist: data.ar?.map(a => a.name).join(', ') || data.artist || '未知歌手',
                cover: data.al?.picUrl || data.cover || '',
                url: data.url || ''
            };
            this.playlist = [song];
            this.updatePlaylistDisplay();
        } catch (error) {
            console.error('加载单曲失败:', error);
            this.showError('加载单曲失败');
        }
    }
    async loadCurrentSong() {
        const song = this.playlist[this.currentIndex];
        if (!song) return;
        try {
            this.elements.songTitle.textContent = song.name || '未知歌曲';
            this.elements.songArtist.textContent = song.artist || '未知歌手';
            if (song.cover) {
                this.elements.albumCover.src = song.cover;
            }
            if (!song.url) {
                const urlResponse = await fetch(`https://api.hypcvgm.top/nmp/v2/url?id=${encodeURIComponent(song.id)}`);
                if (urlResponse.ok) {
                    const urlData = await urlResponse.json();
                    song.url = urlData.url;
                }
            }
            if (song.url) {
                this.audio.src = song.url;
                this.audio.load();
                if (this.config.autoplay) {
                    this.audio.play().catch(() => {});
                }
            }
            await this.loadLyrics(song.id);
            this.updateLyricsDisplay();
            this.updatePlaylistActive();
        } catch (error) {
            console.error('加载当前歌曲失败:', error);
            this.showError('加载歌曲失败');
        }
    }
    async loadLyrics(songId) {
        try {
            const response = await fetch(`https://api.hypcvgm.top/nmp/v2/lyrics?id=${encodeURIComponent(songId)}`);
            if (!response.ok) throw new Error('获取歌词失败');
            const data = await response.json();
            this.lyrics = data.lines || [];
            this.currentLyricIndex = -1;
        } catch (error) {
            console.error('加载歌词失败:', error);
            this.lyrics = [];
            this.currentLyricIndex = -1;
        }
    }
    updateLyricsDisplay() {
        if (this.lyrics.length > 0) {
            this.elements.lyricLine.textContent = this.lyrics[0]?.text || '';
            this.elements.lyricTranslation.textContent = this.lyrics[0]?.translation || '';
        } else {
            this.elements.lyricLine.textContent = '';
            this.elements.lyricTranslation.textContent = '';
        }
    }
    updateCurrentLyric() {
        const currentTime = this.audio.currentTime;
        let index = this.lyrics.findIndex(line => line.time > currentTime);
        if (index === -1) index = this.lyrics.length - 1;
        else index = Math.max(0, index - 1);
        if (index !== this.currentLyricIndex) {
            this.currentLyricIndex = index;
            const line = this.lyrics[index];
            this.elements.lyricLine.textContent = line?.text || '';
            this.elements.lyricTranslation.textContent = line?.translation || '';
            this.elements.lyricsContainer.querySelectorAll('.lyric-line').forEach(el => el.classList.remove('current'));
            this.elements.lyricLine.classList.add('current');
            if (this.elements.lyricTranslation.textContent) {
                this.elements.lyricTranslation.classList.add('current');
            } else {
                this.elements.lyricTranslation.classList.remove('current');
            }
        }
    }
    updateProgressDisplay() {
        this.elements.currentTime.textContent = this.formatTime(this.currentTime);
        this.elements.totalTime.textContent = this.formatTime(this.duration);
        const percent = this.duration > 0 ? (this.currentTime / this.duration) * 100 : 0;
        this.elements.progressBar.style.width = `${percent}%`;
    }
    updateVolumeDisplay() {
        this.elements.volumeBar.style.width = `${(this.audio.volume * 100).toFixed(0)}%`;
    }
    formatTime(seconds) {
        if (!isFinite(seconds) || seconds < 0) return '0:00';
        const mins = Math.floor(seconds / 60);
        const secs = Math.floor(seconds % 60).toString().padStart(2, '0');
        return `${mins}:${secs}`;
    }
    togglePlay() {
        if (this.audio.paused) {
            this.play();
        } else {
            this.pause();
        }
    }
    play() {
        this.audio.play().then(() => {
            this.elements.playIcon.style.display = 'none';
            this.elements.pauseIcon.style.display = 'inline';
        }).catch(() => {});
    }
    pause() {
        this.audio.pause();
        this.elements.playIcon.style.display = 'inline';
        this.elements.pauseIcon.style.display = 'none';
    }
    previousSong() {
        if (this.playlist.length === 0) return;
        if (this.playMode === 'shuffle') {
            this.currentIndex = this.getRandomIndex();
        } else {
            this.currentIndex = (this.currentIndex - 1 + this.playlist.length) % this.playlist.length;
        }
        this.loadCurrentSong();
    }
    nextSong() {
        if (this.playlist.length === 0) return;
        if (this.playMode === 'shuffle') {
            this.currentIndex = this.getRandomIndex();
        } else {
            this.currentIndex = (this.currentIndex + 1) % this.playlist.length;
        }
        this.loadCurrentSong();
    }
    getRandomIndex() {
        let randomIndex;
        do {
            randomIndex = Math.floor(Math.random() * this.playlist.length);
        } while (randomIndex === this.currentIndex);
        return randomIndex;
    }
    handleSongEnd() {
        if (this.playMode === 'single') {
            this.audio.currentTime = 0;
            this.play();
        } else if (this.playMode === 'shuffle') {
            this.currentIndex = this.getRandomIndex();
            this.loadCurrentSong();
        } else {
            this.nextSong();
        }
    }
    updatePlaylistDisplay() {
        if (!this.elements.playlistContent) return;
        this.elements.playlistContent.innerHTML = '';
        if (this.playlist.length === 0) {
            const loadingEl = document.createElement('div');
            loadingEl.className = 'loading';
            loadingEl.textContent = '加载中...';
            this.elements.playlistContent.appendChild(loadingEl);
            return;
        }
        this.playlist.forEach((song, index) => {
            const item = document.createElement('div');
            item.className = 'playlist-item';
            item.dataset.index = index;
            item.innerHTML = `
                <div class="playlist-item-index">${index + 1}</div>
                ${song.cover ? `<img class="playlist-item-cover" src="${song.cover}" alt="封面">` : ''}
                <div class="playlist-item-info">
                    <div class="playlist-item-name">${song.name || '未知歌曲'}</div>
                    <div class="playlist-item-artist">${song.artist || '未知歌手'}</div>
                </div>
            `;
            item.addEventListener('click', () => {
                this.currentIndex = index;
                this.loadCurrentSong();
                this.togglePlaylist(false);
            });
            this.elements.playlistContent.appendChild(item);
        });
        this.updatePlaylistActive();
    }
    updatePlaylistActive() {
        const items = this.elements.playlistContent?.querySelectorAll('.playlist-item') || [];
        items.forEach(item => item.classList.remove('active'));
        const activeItem = Array.from(items).find(item => Number(item.dataset.index) === this.currentIndex);
        if (activeItem) {
            activeItem.classList.add('active');
            activeItem.scrollIntoView({ behavior: 'smooth', block: 'nearest' });
        }
    }
    seekTo(e) {
        if (!this.elements.progressContainer || !this.audio) return;
        const rect = this.elements.progressContainer.getBoundingClientRect();
        const percent = Math.max(0, Math.min(1, (e.clientX - rect.left) / rect.width));
        const newTime = percent * this.duration;
        if (isFinite(newTime) && newTime >= 0) {
            this.audio.currentTime = newTime;
        }
    }
    setVolume(e) {
        if (!this.elements.volumeSlider) return;
        const rect = this.elements.volumeSlider.getBoundingClientRect();
        const percent = Math.max(0, Math.min(1, (e.clientX - rect.left) / rect.width));
        this.volume = percent;
        this.audio.volume = this.volume;
        this.updateVolumeDisplay();
    }
    toggleLyrics() {
        this.showLyrics = !this.showLyrics;
        this.elements.lyricsContainer.classList.toggle('hidden', !this.showLyrics);
        this.elements.lyricsBtn.classList.toggle('active', this.showLyrics);
    }
    togglePlaylist(show = null) {
        if (!this.elements.playlistContainer) return;
        const isShowing = this.elements.playlistContainer.classList.contains('show');
        const shouldShow = show !== null ? show : !isShowing;
        if (shouldShow) {
            this.determinePlaylistDirection();
            this.updatePlaylistDisplay();
            this.elements.playlistContainer.classList.add('show');
            if (this.elements.listBtn) {
                this.elements.listBtn.classList.add('active');
            }
        } else {
            this.elements.playlistContainer.classList.remove('show', 'show-above', 'show-below');
            if (this.elements.listBtn) {
                this.elements.listBtn.classList.remove('active');
            }
        }
    }
    togglePlayMode() {
        const modes = ['list', 'single', 'shuffle'];
        const currentIndex = modes.indexOf(this.playMode);
        this.playMode = modes[(currentIndex + 1) % 3];
        
        const iconSvgs = { list: ICONS.loopList, single: ICONS.loopSingle, shuffle: ICONS.shuffle };
        const titles = { list: '列表循环', single: '单曲循环', shuffle: '随机播放' };
        
        if (this.elements.loopModeBtn) {
            this.elements.loopModeBtn.innerHTML = iconSvgs[this.playMode];
            this.elements.loopModeBtn.title = titles[this.playMode];
        }
    }
    toggleMinimize() {
        const isCurrentlyMinimized = this.element.classList.contains('minimized');
        this.isMinimized = isCurrentlyMinimized;
        if (!isCurrentlyMinimized) {
            this.element.classList.add('minimized');
            this.isMinimized = true;
            if (this.elements.minimizeBtn) {
                this.elements.minimizeBtn.classList.add('active');
                this.elements.minimizeBtn.title = '展开';
                this.elements.minimizeBtn.innerHTML = ICONS.maximize; 
            }
            this.clearIdleTimer();
            this.isIdle = false;
            this.element.classList.remove('idle', 'fading-in', 'fading-out', 'docked-left', 'docked-right', 'popping-left', 'popping-right');
            this.startIdleTimer();
        } else {
            this.element.classList.remove('minimized');
            this.isMinimized = false;
            if (this.elements.minimizeBtn) {
                this.elements.minimizeBtn.classList.remove('active');
                this.elements.minimizeBtn.title = '缩小';
                this.elements.minimizeBtn.innerHTML = ICONS.minimize;
            }
            this.clearIdleTimer();
            if (this.isIdle) {
                this.restoreOpacity();
            } else {
                this.element.classList.remove('idle', 'fading-in', 'fading-out', 'docked-left', 'docked-right', 'popping-left', 'popping-right');
            }
            this.isIdle = false;
        }
    }
    determinePlaylistDirection() {
        const playerRect = this.element.getBoundingClientRect();
        const viewportHeight = window.innerHeight;
        const spaceBelow = viewportHeight - playerRect.bottom;
        const spaceAbove = playerRect.top;
        const playlistHeight = 220;
        this.elements.playlistContainer.classList.remove('expand-up');
        if (spaceBelow >= playlistHeight || spaceBelow >= spaceAbove) {
        } else {
            this.elements.playlistContainer.classList.add('expand-up');
        }
    }
    setupDragAndDrop() {
        return;
    }
    showError(message) {
        this.elements.songTitle.textContent = message;
        this.elements.songArtist.textContent = '';
        this.elements.lyricLine.textContent = '';
    }
    initTheme() {
        this.setTheme(this.config.theme);
        if (this.config.theme === 'auto') {
            this.setupThemeListener();
        }
    }
    setTheme(theme) {
        if (theme === 'auto') {
            const detectedTheme = this.detectTheme();
            this.element.setAttribute('data-theme', 'auto');
            if (detectedTheme === 'dark') {
                this.element.classList.add('theme-dark-detected');
            } else {
                this.element.classList.remove('theme-dark-detected');
            }
        } else {
            this.element.setAttribute('data-theme', theme);
            this.element.classList.remove('theme-dark-detected');
        }
    }
    detectTheme() {
        const hostTheme = this.detectHostTheme();
        if (hostTheme) {
            return hostTheme;
        }
        const cssTheme = this.detectCSSTheme();
        if (cssTheme) {
            return cssTheme;
        }
        return this.detectSystemTheme();
    }
    detectHostTheme() {
        const html = document.documentElement;
        const body = document.body;
        const darkClasses = ['dark', 'theme-dark', 'dark-theme', 'dark-mode'];
        const lightClasses = ['light', 'theme-light', 'light-theme', 'light-mode'];
        for (const className of darkClasses) {
            if (html.classList.contains(className)) return 'dark';
        }
        for (const className of lightClasses) {
            if (html.classList.contains(className)) return 'light';
        }
        if (body) {
            for (const className of darkClasses) {
                if (body.classList.contains(className)) return 'dark';
            }
            for (const className of lightClasses) {
                if (body.classList.contains(className)) return 'light';
            }
        }
        const htmlTheme = html.getAttribute('data-theme');
        if (htmlTheme === 'dark' || htmlTheme === 'light') {
            return htmlTheme;
        }
        const bodyTheme = body?.getAttribute('data-theme');
        if (bodyTheme === 'dark' || bodyTheme === 'light') {
            return bodyTheme;
        }
        return null;
    }
    detectCSSTheme() {
        try {
            const rootStyles = getComputedStyle(document.documentElement);
            const bgColor = rootStyles.getPropertyValue('--bg-color') || 
                           rootStyles.getPropertyValue('--background-color') ||
                           rootStyles.getPropertyValue('--color-bg');
            const textColor = rootStyles.getPropertyValue('--text-color') || 
                             rootStyles.getPropertyValue('--color-text') ||
                             rootStyles.getPropertyValue('--text-primary');
            if (bgColor || textColor) {
                const isDarkBg = this.isColorDark(bgColor);
                const isLightText = this.isColorLight(textColor);
                if (isDarkBg || isLightText) {
                    return 'dark';
                }
                if (!isDarkBg || !isLightText) {
                    return 'light';
                }
            }
        } catch (error) {
            console.warn('CSS主题检测失败:', error);
        }
        return null;
    }
    detectSystemTheme() {
        if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
            return 'dark';
        }
        return 'light';
    }
    isColorDark(color) {
        if (!color) return false;
        color = color.replace(/\s/g, '').toLowerCase();
        if (color.includes('dark') || color.includes('black') || color === 'transparent') {
            return true;
        }
        const rgb = color.match(/rgb\((\d+),(\d+),(\d+)\)/);
        if (rgb) {
            const [, r, g, b] = rgb.map(Number);
            const brightness = (r * 299 + g * 587 + b * 114) / 1000;
            return brightness < 128;
        }
        const hex = color.match(/^#([0-9a-f]{3}|[0-9a-f]{6})$/);
        if (hex) {
            const hexValue = hex[1];
            const r = parseInt(hexValue.length === 3 ? hexValue[0] + hexValue[0] : hexValue.substr(0, 2), 16);
            const g = parseInt(hexValue.length === 3 ? hexValue[1] + hexValue[1] : hexValue.substr(2, 2), 16);
            const b = parseInt(hexValue.length === 3 ? hexValue[2] + hexValue[2] : hexValue.substr(4, 2), 16);
            const brightness = (r * 299 + g * 587 + b * 114) / 1000;
            return brightness < 128;
        }
        return false;
    }
    isColorLight(color) {
        return !this.isColorDark(color);
    }
    setupThemeListener() {
        if (window.matchMedia) {
            const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)');
            const handleThemeChange = () => {
                if (this.config.theme === 'auto') {
                    this.setTheme('auto');
                }
            };
            if (mediaQuery.addEventListener) {
                mediaQuery.addEventListener('change', handleThemeChange);
            } else {
                mediaQuery.addListener(handleThemeChange);
            }
        }
        if (window.MutationObserver) {
            const observer = new MutationObserver((mutations) => {
                if (this.config.theme === 'auto') {
                    let shouldUpdate = false;
                    mutations.forEach((mutation) => {
                        if (mutation.type === 'attributes' && 
                            (mutation.attributeName === 'class' || mutation.attributeName === 'data-theme')) {
                            shouldUpdate = true;
                        }
                    });
                    if (shouldUpdate) {
                        this.setTheme('auto');
                    }
                }
            });
            observer.observe(document.documentElement, {
                attributes: true,
                attributeFilter: ['class', 'data-theme']
            });
            if (document.body) {
                observer.observe(document.body, {
                    attributes: true,
                    attributeFilter: ['class', 'data-theme']
                });
            }
        }
    }
    static init() {
        document.querySelectorAll('.netease-mini-player').forEach(element => {
            if (!element._neteasePlayer) {
                element._neteasePlayer = new NeteaseMiniPlayer(element);
            }
        });
    }
    static initPlayer(element) {
        if (!element._neteasePlayer) {
            element._neteasePlayer = new NeteaseMiniPlayer(element);
        }
        return element._neteasePlayer;
    }
}
if (typeof window !== 'undefined') {
    window.NeteaseMiniPlayer = NeteaseMiniPlayer;
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', NeteaseMiniPlayer.init);
    } else {
        NeteaseMiniPlayer.init();
    }
}

class NMPv2ShortcodeParser {
    constructor() {
        this.paramMappings = {
            'position': 'data-position',
            'theme': 'data-theme',
            'lyric': 'data-lyric',
            'embed': 'data-embed',
            'minimized': 'data-default-minimized',
            'autoplay': 'data-autoplay',
            'idle-opacity': 'data-idle-opacity',
            'auto-pause': 'data-auto-pause'
        };

        if (document.readyState === 'loading') {
            document.addEventListener('DOMContentLoaded', () => this.init());
        } else {
            this.init();
        }
    }

    init() {
        this.processContainer(document.body);
    }

    /**
     * 处理容器内的所有短语法
     */
    processContainer(container) {
        this.processTextNodes(container);
        this.processExistingElements(container);
        this.initializePlayers(container);
    }

    /**
     * 处理文本节点中的短语法
     */
    processTextNodes(container) {
        const walker = document.createTreeWalker(
            container, 
            NodeFilter.SHOW_TEXT, 
            null, 
            false
        );

        const textNodes = [];
        let node;
        while (node = walker.nextNode()) {
            if (node.textContent.includes('{nmpv2:')) {
                textNodes.push(node);
            }
        }

        textNodes.forEach(node => {
            const content = node.textContent;
            const shortcodes = this.extractShortcodes(content);
            
            if (shortcodes.length > 0) {
                const fragment = document.createDocumentFragment();
                let lastIndex = 0;

                shortcodes.forEach(shortcode => {
                    if (shortcode.startIndex > lastIndex) {
                        fragment.appendChild(document.createTextNode(
                            content.substring(lastIndex, shortcode.startIndex)
                        ));
                    }

                    const playerElement = this.createPlayerElement(shortcode);
                    fragment.appendChild(playerElement);

                    lastIndex = shortcode.endIndex;
                });

                if (lastIndex < content.length) {
                    fragment.appendChild(document.createTextNode(
                        content.substring(lastIndex)
                    ));
                }

                node.parentNode.replaceChild(fragment, node);
            }
        });
    }

    processExistingElements(container) {
        container.querySelectorAll('.netease-mini-player:not([data-shortcode-processed])')
            .forEach(element => {
                element.setAttribute('data-shortcode-processed', 'true');
            });
    }

    initializePlayers(container) {
        container.querySelectorAll('.netease-mini-player:not([data-initialized])')
            .forEach(element => {
                element.setAttribute('data-initialized', 'true');
                NeteaseMiniPlayer.initPlayer(element);
            });
    }

    extractShortcodes(text) {
        const regex = /\{nmpv2:([^}]*)\}/g;
        let match;
        const results = [];
        let lastIndex = 0;

        while ((match = regex.exec(text)) !== null) {
            const content = match[1].trim();
            const startIndex = match.index;
            const endIndex = match.index + match[0].length;

            let shortcode = {
                type: 'song',
                id: null,
                params: {},
                startIndex,
                endIndex
            };

            this.parseShortcodeContent(content, shortcode);
            results.push(shortcode);
        }

        return results;
    }

    parseShortcodeContent(content, shortcode) {
        if (content.startsWith('playlist=')) {
            shortcode.type = 'playlist';
            const parts = content.split(/,\s*/);
            const firstPart = parts.shift();
            shortcode.id = firstPart.replace('playlist=', '').trim();
            
            parts.forEach(part => this.parseParam(part, shortcode.params));
        } else if (content.includes('=')) {
            const parts = content.split(/,\s*/);
            const firstPart = parts.shift();
            
            if (firstPart.includes('=')) {
                this.parseParam(firstPart, shortcode.params);
                parts.forEach(part => this.parseParam(part, shortcode.params));
            } else {
                shortcode.id = firstPart.trim();
                parts.forEach(part => this.parseParam(part, shortcode.params));
            }
        } else {
            shortcode.id = content.trim();
        }

        if (shortcode.params.position === undefined || shortcode.params.position === 'static') {
            shortcode.params.embed = shortcode.params.embed ?? 'true';
        } else if (shortcode.params.embed === undefined) {
            shortcode.params.embed = 'false';
        }
    }

    parseParam(paramStr, params) {
        const [key, value] = paramStr.split('=');
        if (!key || !value) return;

        const cleanKey = key.trim().toLowerCase();
        const cleanValue = value.trim().toLowerCase();

        if (cleanKey === 'song-id') {
            params.songId = cleanValue;
        } else if (cleanKey === 'playlist-id') {
            params.playlistId = cleanValue;
            params.type = 'playlist';
        } else if (cleanKey === 'minimized') {
            params.defaultMinimized = cleanValue === 'true' ? 'true' : 'false';
        } else {
            const mapping = this.paramMappings[cleanKey] || `data-${cleanKey}`;
            params[cleanKey] = cleanValue;
        }
    }

    createPlayerElement(shortcode) {
        const div = document.createElement('div');
        div.className = 'netease-mini-player';
        div.setAttribute('data-shortcode-processed', 'true');

        if (shortcode.type === 'playlist' && shortcode.id) {
            div.setAttribute('data-playlist-id', shortcode.id);
        } else if (shortcode.id) {
            div.setAttribute('data-song-id', shortcode.id);
        }

        Object.entries(shortcode.params).forEach(([key, value]) => {
            if (key === 'songId') {
                div.setAttribute('data-song-id', value);
            } else if (key === 'playlistId') {
                div.setAttribute('data-playlist-id', value);
            } else if (key === 'type') {
            } else {
                const dataKey = this.paramMappings[key] || `data-${key}`;
                div.setAttribute(dataKey, value);
            }
        });

        return div;
    }

    static processDynamicContent(content) {
        const tempDiv = document.createElement('div');
        tempDiv.innerHTML = content;
        window.nmpv2ShortcodeParser.processContainer(tempDiv);
        return tempDiv.innerHTML;
    }
}

if (typeof window !== 'undefined') {
    window.nmpv2ShortcodeParser = new NMPv2ShortcodeParser();
    
    window.processNMPv2Shortcodes = function(container) {
        if (container instanceof Element) {
            window.nmpv2ShortcodeParser.processContainer(container);
        } else {
            console.warn('processNMPv2Shortcodes requires a DOM element');
        }
    };
}

if (typeof module !== 'undefined' && module.exports) {
    module.exports = {
        renderShortcodes: function(html) {
            return html.replace(/\{nmpv2:([^}]*)\}/g, (match, content) => {
                let shortcode = {
                    type: 'song',
                    id: null,
                    params: {}
                };
                
                if (content.startsWith('playlist=')) {
                    shortcode.type = 'playlist';
                    const parts = content.split(/,\s*/);
                    shortcode.id = parts[0].replace('playlist=', '').trim();
                    parts.slice(1).forEach(part => {
                        const [key, value] = part.split('=');
                        if (key && value) shortcode.params[key.trim()] = value.trim();
                    });
                } else {
                    const parts = content.split(/,\s*/);
                    if (parts[0].includes('=')) {
                        parts.forEach(part => {
                            const [key, value] = part.split('=');
                            if (key && value) shortcode.params[key.trim()] = value.trim();
                        });
                    } else {
                        shortcode.id = parts[0].trim();
                        parts.slice(1).forEach(part => {
                            const [key, value] = part.split('=');
                            if (key && value) shortcode.params[key.trim()] = value.trim();
                        });
                    }
                }

                if (!shortcode.params.position || shortcode.params.position === 'static') {
                    shortcode.params.embed = shortcode.params.embed ?? 'true';
                } else if (shortcode.params.embed === undefined) {
                    shortcode.params.embed = 'false';
                }

                let html = '<div class="netease-mini-player"';
                
                if (shortcode.type === 'playlist' && shortcode.id) {
                    html += ` data-playlist-id="${shortcode.id}"`;
                } else if (shortcode.id) {
                    html += ` data-song-id="${shortcode.id}"`;
                }

                Object.entries(shortcode.params).forEach(([key, value]) => {
                    if (key === 'songId') {
                        html += ` data-song-id="${value}"`;
                    } else if (key === 'playlistId') {
                        html += ` data-playlist-id="${value}"`;
                    } else {
                        const dataKey = {
                            'position': 'data-position',
                            'theme': 'data-theme',
                            'lyric': 'data-lyric',
                            'embed': 'data-embed',
                            'minimized': 'data-default-minimized',
                            'autoplay': 'data-autoplay',
                            'idle-opacity': 'data-idle-opacity',
                            'auto-pause': 'data-auto-pause'
                        }[key] || `data-${key}`;
                        html += ` ${dataKey}="${value}"`;
                    }
                });

                html += '></div>';
                return html;
            });
        }
    };
}

console.log(["版本号 v2.1.0", "NeteaseMiniPlayer V2 [NMPv2]", "BHCN STUDIO & 北海的佰川（ImBHCN[numakkiyu]）", "GitHub地址：https://github.com/numakkiyu/NeteaseMiniPlayer", "基于 Apache 2.0 开源协议发布"].join("\n"));
