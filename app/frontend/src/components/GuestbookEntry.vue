<script setup lang="ts">
    import { computed } from 'vue'
    import { StarIcon }  from '@heroicons/vue/24/solid'
    import { StarIcon as HollowStarIcon} from '@heroicons/vue/24/outline'

    const pictureSource = computed(() => {
        return  props.github ? 'https://github.com/' + props.github + '.png' : '/img/dummy.png'
    })

    const userName = computed(() => {
        return props.lastName ? props.lastName.concat(', ', props.firstName) : props.firstName 
    })

    const props = defineProps<{
        firstName: string,
        lastName: string,
        occupation: string,
        github: string,
        rating: number,
        testimonial: String,
    }>()
</script>

<template>
    <div class="guestbook-entry">
        <div class="user-info">
            <img class="profile-picture" :src="pictureSource" />
            <div class="name-and-job">
                <b>{{ userName }} </b>
                <br/>
                <span class="occupation">{{ occupation ? occupation : '' }}</span>
            </div>
            <div class="rating">
                <StarIcon v-for="i in rating" :key="i" class="icon"/>
                <HollowStarIcon v-for="i in 5 - rating" :key="i" class="icon"/>
            </div>
        </div>
        <div class="testimonial">
            <i>{{ testimonial }}</i>
        </div>
    </div>
</template>

<style scoped>
    .guestbook-entry {
        @apply border-light;
        @apply border-dashed;
        @apply border-2;
        @apply rounded-xl;

        @apply p-4;
        @apply m-4;
        @apply max-w-md;

    }

    .profile-picture {
        @apply object-contain;
        @apply w-16;
        @apply h-16;
        @apply rounded-full;
    }

    .user-info {
        @apply flex;
        @apply flex-row;
    }

    .name-and-job {
        @apply ml-4;   
        @apply text-brand;
    }

    .occupation {
        @apply text-mid;
    }

    .testimonial {
        @apply text-light;
        @apply mt-4;
    }

    .rating {
        @apply flex;
        @apply justify-end;
        @apply grow;
    }

    
</style>