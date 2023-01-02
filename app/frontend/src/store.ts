import { defineStore } from 'pinia'

export const useSortOrderStore = defineStore('sortOrder', {
    state: () => {
        return { state: 'newest' }
    },
    getters: {
        sortOrder: (state) => state.state
    },
    actions: {
        change() {
            if (this.state == 'newest') {
                this.state = 'best'
            }
            else if (this.state == 'best') {
                this.state = 'worst'
            } else {
                this.state = 'newest'
            }
        },
    },
})

export const useShowModalStore = defineStore('modalState', {
    state: () => {
        return { state: false}
    },
    getters: {
        modalState: (state) => state.state
    },
    actions: {
        invert() {
            this.state = !this.state
        }
    },
})