<script setup lang="ts">
import { useForm } from "@tanstack/vue-form";
import { useRouter } from "vue-router";
import useAuth from "@/shared/hooks/auth/useAuth";
import { useRegister } from "@/features/auth/queries";
import {
	registerFormSchema,
	type RegisterFormData,
} from "@/features/auth/schemas";
import { Box, Button, InputField, Text } from "@/shared/components/ui";
import { ArrowPath } from "@/shared/components/icons";

const router = useRouter();
const { setUser } = useAuth();
const { mutateAsync: register, isPending } = useRegister();

const form = useForm({
	defaultValues: {
		email: "",
		password: "",
		confirmPassword: "",
	},
	validators: {
		onChange: registerFormSchema,
	},
	onSubmit: async ({ value }) => {
		const user = await register(value);
		setUser(user);
		router.push({ name: "Home" });
	},
});
</script>

<template>
	<Box bg="darkest" width="w-full" height="h-full" class="items-center">
		<Text as="h3" size="3xl">Create your account </Text>
		<form
			@submit="
        (e) => {
          e.preventDefault();
          e.stopPropagation();
          form.handleSubmit();
        }
      "
			class="w-4/5 sm:w-2/5 flex flex-col gap-4"
		>
			<!-- Email Field -->
			<form.Field name="email">
				<template v-slot="{ field }">
					<InputField
						:id="field.name"
						:name="field.name"
						:value="field.state.value"
						bg="primary"
						text-color="dark"
						type="email"
						:disabled="isPending"
						errorslot
						@input="(value: string | number | null) => {
              if (typeof value !== 'string') return;
              field.handleChange(value);
            }"
						@blur="field.handleBlur"
					>
						<template #label>
							<Text>Email</Text>
						</template>
						<template
							#error
							v-if="field.state.meta.isTouched && field.state.meta.errors.length > 0"
						>
							<Text color="error">
								{{ field.state.meta.errors[0]?.message }}
							</Text>
						</template>
					</InputField>
				</template>
			</form.Field>

			<!-- Password Field -->
			<form.Field name="password">
				<template v-slot="{ field }">
					<InputField
						:id="field.name"
						:name="field.name"
						:value="field.state.value"
						text-color="dark"
						bg="primary"
						type="password"
						:disabled="isPending"
						errorslot
						@input="(value: string | number | null) => {
              if (typeof value !== 'string') return;
              field.handleChange(value);
            }"
						@blur="field.handleBlur"
					>
						<template #label>
							<Text>Password</Text>
						</template>
						<template
							#error
							v-if="field.state.meta.isTouched && field.state.meta.errors.length > 0"
						>
							<Text color="error">
								{{ field.state.meta.errors[0]?.message }}
							</Text>
						</template>
					</InputField>
				</template>
			</form.Field>

			<!-- Confirm Password Field -->
			<form.Field name="confirmPassword">
				<template v-slot="{ field }">
					<InputField
						:id="field.name"
						:name="field.name"
						:value="field.state.value"
						text-color="dark"
						bg="primary"
						type="password"
						:disabled="isPending"
						errorslot
						@input="(value: string | number | null) => {
              if (typeof value !== 'string') return;
              field.handleChange(value);
            }"
						@blur="field.handleBlur"
					>
						<template #label>
							<Text>Confirm Password</Text>
						</template>
						<template
							#error
							v-if="field.state.meta.isTouched && field.state.meta.errors.length > 0"
						>
							<Text color="error">
								{{ field.state.meta.errors[0]?.message }}
							</Text>
						</template>
					</InputField>
				</template>
			</form.Field>

			<!-- Submit Button -->
			<form.Subscribe>
				<template v-slot="{ canSubmit, isSubmitting }">
					<Button
						type="submit"
						class="mt-4"
						height="h-10"
						width="w-full"
						:disabled="!canSubmit || isPending || isSubmitting"
					>
						<ArrowPath class="animate-spin" v-if="isSubmitting || isPending"/>
						<Text v-else>Register</Text>
					</Button>
				</template>
			</form.Subscribe>
		</form>
	</Box>
</template>
