<script setup lang="ts">
    import { XMarkIcon } from '@heroicons/vue/24/solid';
    import axios from 'axios';
    import { ref } from 'vue';
    import RatingPicker from './RatingPicker.vue';
    import { useShowModalStore } from '@/store'

    const showModalStore = useShowModalStore()

    var warning = ref(false)

    const firstname = ref<string | undefined>()
    const lastname = ref<string | undefined>()
    const occupation = ref<string | undefined>()
    const github = ref<string | undefined>()
    const testimonial = ref<string | undefined>()
    const rating = ref(3)

    const emit = defineEmits<{
        (e: 'createdEntry', result: [notification: string, success: boolean] ): void
    }>()

    function isValidString(s: string | undefined): Boolean {
        return  s != undefined && s != ""
    }

    function isValidRating(r: number): Boolean {
        return r > 0 && r < 6
    }

    function postTestimonial(){
        if (isValidString(firstname.value) && isValidString(testimonial.value) && isValidRating(rating.value)) {
            warning.value = false

            axios.post("/v1/entries", {
                firstname: firstname.value,
                lastname: lastname.value ? lastname.value : "",
                occupation: occupation.value ? occupation.value : "",
                github: github.value ? github.value : "",
                testimonial: testimonial.value,
                rating: rating.value
            }).then(function(response) {
                    showModalStore.invert();
                    if (response.data.status == "failure") {
                        emit('createdEntry', ["You already submitted feedback!", false])
                    } else {
                        emit('createdEntry', ["Your feedback was submitted successfully.", true])
                    }
                }
            ).catch(e => {
                showModalStore.invert()
                emit('createdEntry', ["An error occured while submitting your feedback. Please try again.", false])
            })
        } else {
            warning.value = true
        }
    }

</script>

<template>
    <!--modal content-->
    <div class="overlay" v-if="showModalStore.state">
        <div class="modalwindow">
            <div class="close-button">
                <XMarkIcon @click="showModalStore.invert()" class="icon clickable"/>
            </div>
            <div v-if="warning" class="warning">
                Fields marked with '*' are required!
            </div>
            <form action="/v1/entries" method="POST">
                <div>
                    <div class="form-line">
                        <label for="firstname">First Name*</label>
                        <input v-model="firstname" id="firstname" name="firstname" type="text" placeholder="First Name" required />
                    </div>
                    <div class="form-line">
                        <label for="lastname">Last Name</label>
                        <input v-model="lastname" id="lastname" name="lastname" type="text" placeholder="Last Name"/>
                    </div>
                    <div class="form-line">
                        <label for="occupation">Occupation</label>
                        <input v-model="occupation" id="occupation" name="occupation" type="text" placeholder="Occupation"/>
                    </div>
                    <div class="form-line">
                        <label for="github">GitHub Handle</label>
                        <input v-model="github" id="github" name="github" type="text" placeholder="mocdaniel"/>
                    </div>
                    <div class="form-line">
                        <label for="rating">Rating</label>
                        <RatingPicker @setRating="(newRating) => rating = newRating"  :rating="rating"/>
                    </div>
                    <div class="form-line">
                        <label for="feedback">Feedback*</label>
                        <textarea v-model="testimonial" id="testimonial" name="testimonial" rows="4" cols="25" required/> 
                    </div>
                    <div class="button-line">
                        <button class="blue-button" @click="postTestimonial" type="button" value="Submit">Submit</button>
                    </div>
                </div>
            </form>
        </div>
    </div>
</template>

<style scoped>
    textarea {
        @apply ml-4;
    }
    .button-line {
        @apply flex;
        @apply flex-row;
        @apply justify-end;
    }

    
    .close-button {
        @apply flex;
        @apply flex-row;
        @apply justify-end;
        @apply mb-4;
    }

    .form-line {
        @apply flex;
        @apply flex-row;
        @apply justify-between;
        @apply mb-2;
    }
    .icon {
        @apply text-brand;
        @apply stroke-brand;
    }
    .modalwindow {
        @apply relative;
        @apply top-1/3;
        @apply mx-auto;
        @apply p-5;
        @apply w-fit;
        @apply shadow-lg;
        @apply rounded-md;
        @apply bg-light;
        @apply border-brand;
        @apply border;
    }

    .overlay {
        @apply fixed;
        @apply inset-0;
        @apply overflow-y-auto;
        @apply bg-transparent;
        @apply h-full;
        @apply w-full;
        @apply backdrop-filter;
        @apply backdrop-blur-sm;
    }

    .warning {
        @apply px-2;
        @apply py-1;
        @apply mb-4;
        @apply bg-red-400;
        @apply text-light;
    }
</style>