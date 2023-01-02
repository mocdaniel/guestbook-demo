<script setup lang="ts">
    import {computed, ref, type Ref } from 'vue'
    import GuestbookEntry from './GuestbookEntry.vue'
    import { useSortOrderStore } from '@/store'
    import axios from 'axios'

    const sortOrderStore = useSortOrderStore()
    const data: Ref<Entries> = ref({entries: []})

    type Entry = {
        id: number,
        firstname: string,
        lastname: string,
        occupation: string,
        github: string,
        rating: number,
        testimonial: string,
    }

    type Entries = {
        entries: Entry[]
    }

    function loadEntries() {
        axios.get<Entries>('/v1/entries').then(function(response) {
            data.value = response.data
            data.value.entries.sort((a: { id: number; }, b: { id: number; }) => b.id - a.id)
        })
    }

    const getAvgRating = computed(() => {
        if (data.value.entries.length > 0) {
            var addedRatings = data.value.entries.reduce((acc, obj) => {
                return acc + obj.rating
            }, 0)
            return addedRatings / data.value.entries.length
        } else {
            return 0.0
        }
    })

    defineExpose({
        loadEntries,
        getAvgRating,
    })

    sortOrderStore.$subscribe((mutation, newSortOrder) => {
        if (newSortOrder.state == 'newest') {
            data.value.entries.sort((a: { id: number; }, b: { id: number; }) => b.id - a.id)
        } else if (newSortOrder.state == 'worst') {
            data.value.entries.sort((a: { rating: number; }, b: { rating: number; }) => a.rating - b.rating)
        } else {
            data.value.entries.sort((a: { rating: number; }, b: { rating: number; }) => b.rating - a.rating)
        }
    })

    loadEntries()
    setInterval(function() {
        loadEntries()
    }.bind(this), 30000)
</script>

<template>
    <div class="guestbook-entries">
        <GuestbookEntry
            v-for="entry in data.entries"
            :key="entry.id"
            :entry="entry"
            :firstName="entry.firstname"
            :lastName="entry.lastname"
            :occupation="entry.occupation"
            :github="entry.github"
            :rating="entry.rating"
            :testimonial="entry.testimonial"
        />
    </div>
</template>

<style scoped>
    .guestbook-entries {
        @apply flex-col;
        @apply flex-wrap;
        @apply justify-center;
        @apply h-full;
        @apply w-1/2;
        @apply mx-auto;
        @apply max-w-lg;
        @apply pt-8;
    }
</style>