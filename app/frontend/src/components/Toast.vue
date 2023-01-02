<script setup lang="ts">
    import { CheckCircleIcon, ExclamationCircleIcon } from '@heroicons/vue/24/outline';

    const props = defineProps<{
        notification: string | null,
        success: boolean,
    }>()

</script>

<template>
    <div v-if="notification" class="toast" :class="success ? 'green' : 'red'">
        <div class="icon">
            <CheckCircleIcon v-if="success" />
            <ExclamationCircleIcon v-else/>
        </div>
        <div>
            {{ notification }}
        </div>
    </div>
</template>

<style scoped>
    .toast {
        @apply fixed;
        @apply right-4;
        @apply top-4;
        @apply h-fit;
        @apply w-fit;
        @apply m-4;
        @apply p-2;
        @apply text-light;
        @apply rounded-md;
        @apply flex;
        @apply flex-row;
        @apply justify-between;
    }

    .red {
        @apply bg-red-400
    }

    .green {
        @apply bg-green-500;
    }

    .icon {
        @apply text-light;
        @apply stroke-light;
    }
</style>