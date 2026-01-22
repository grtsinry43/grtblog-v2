<script lang="ts">
</script>

<!-- Ambient Background Blurs (Global) -->
<div class="fixed top-0 left-0 w-full h-screen pointer-events-none z-0 overflow-hidden">
    <!-- Top Left Blob: Gold/Stone hue in dark mode, subtle stone in light -->
    <div class="absolute top-[-10%] left-[-10%] w-[60vw] h-[60vw] bg-stone-200/40 dark:bg-amber-900/10 rounded-full blur-[100px] mix-blend-multiply dark:mix-blend-screen transition-colors duration-1000"></div>
    
    <!-- Bottom Right Blob: Jade/Moss hue -->
    <div class="absolute bottom-[-10%] right-[-10%] w-[50vw] h-[50vw] bg-jade-100/40 dark:bg-jade-900/15 rounded-full blur-[120px] mix-blend-multiply dark:mix-blend-screen transition-colors duration-1000"></div>
</div>

<!-- Noise Texture -->
<div class="bg-noise"></div>
