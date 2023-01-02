<script setup lang="ts">
    import { computed, ref } from 'vue'
    import { StarIcon }  from '@heroicons/vue/24/solid'
    import { StarIcon as HollowStarIcon} from '@heroicons/vue/24/outline'

    const currentRating = computed(() => {
        return props.rating
    })

    defineEmits<{
        (e: 'setRating', id: number): void
    }>()

    const props = defineProps<{
        rating: number,
    }>()
</script>

<template>
    <div class="star-row">
        <StarIcon @click="$emit('setRating', i)" v-for="i in currentRating" class="icon"/>
        <HollowStarIcon @click="$emit('setRating', i + rating)" v-for="i in 5 - currentRating" class="icon"/>
        <input v-model="rating" id="rating" name="rating" type="number" min="1" max="5" hidden/>
    </div>
</template>

<style scoped>
    .star-row {
        @apply flex;
        @apply flex-row;
        @apply justify-end;
        @apply cursor-pointer;
    }
</style>