<script setup lang="ts">
    import Guestbook from './components/Guestbook.vue'
    import Controls from './components/Controls.vue'
    import FeedbackModal from './components/FeedbackModal.vue';
    import Toast from './components/Toast.vue'
    import { ref } from 'vue';

    const notification = ref<string | null>(null)
    const success = ref(false)
    const guestbook = ref<InstanceType<typeof Guestbook> | null>(null)

    const delay = (ms: number) => new Promise(res => setTimeout(res, ms));

    async function displayNotification(this: any, n: string, s: boolean) {
        notification.value = n
        success.value = s

        if (s && guestbook.value) {
            guestbook.value.loadEntries()
        }

        await delay(10000)

        notification.value = null
    }
</script>

<template>
    <div class="all">
        <header>
            <div>
                <h1 class="heading">Workshop Guestbook</h1>
                <Controls :ratings="guestbook?.getAvgRating"/>
            </div>
        </header>
        <main>
            <div>
                <Guestbook ref="guestbook"/>
            </div>
        </main>
    </div>
    <FeedbackModal @created-entry="(result) => displayNotification(...result)"/>
    <Toast :notification="notification" :success="success"/>
</template>

<style scoped>
    header {
        @apply p-4;
        @apply text-white;
        @apply sticky;
    }

    .all {
        @apply h-screen;
        @apply flex;
        @apply flex-col;
    }

    .logo {
        @apply mx-auto;
    }

    .heading {
        @apply mt-8;
        @apply mx-auto;
        @apply text-center;
    }

    main {
        @apply overflow-scroll;
        @apply flex-grow;
        @apply flex-wrap;

        background-size: 16px 16px;
        background-image:
            linear-gradient(to right, #111828 1px, transparent 1px),
            linear-gradient(to bottom, #111828 1px, transparent 1px);    
    }
</style>
