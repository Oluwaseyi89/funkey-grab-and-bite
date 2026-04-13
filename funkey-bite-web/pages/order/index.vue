<template>
    <div class="mt-[-30px]">
     
      <PageHeader
      title-before="Your"
      highlight-text="Order"
      title-after=""
      subtitle="Review your items and complete your order"
      :narrow="true"
      :alignment="'center'"
      variant="gradient"
    />
  
      <div class="section-padding mx-3 md:mx-8 mt-[-60px] md:mt-[-80px] px-3 py-3 md:py-8 md:px-8">
        <div class="container-narrow">

          <div v-if="cart.items.length === 0" class="text-center py-12">
            <ShoppingCart class="w-24 h-24 text-gray-300 dark:text-gray-700 mx-auto mb-6" />
            <h2 class="text-2xl font-bold text-gray-900 dark:text-white mb-4">Your cart is empty</h2>
            <p class="text-gray-600 dark:text-gray-400 mb-8 max-w-md mx-auto">
              Add some delicious items from our menu to get started!
            </p>
            <NuxtLink to="/menu" class="btn-primary text-lg px-8 py-4">
              <ShoppingBag class="w-5 h-5 inline mr-2" />
              Browse Menu
            </NuxtLink>
          </div>
  

          <div v-else>

            <div class="my-5 md:my-8 px-3 py-3 md:py-8 md:px-8 mt-0">
              <div class="flex items-center justify-center space-x-4 md:space-x-12">
                <div v-for="(step, index) in steps" :key="step.id" class="flex items-center">
                  <div class="flex flex-col items-center">
                    <div
                      class="w-10 h-10 rounded-full flex items-center justify-center font-bold transition-all"
                      :class="[
                        currentStep >= index
                          ? 'bg-brand-500 text-white'
                          : 'bg-gray-200 dark:bg-slate-700 text-gray-500 dark:text-gray-400'
                      ]"
                    >
                      {{ index + 1 }}
                    </div>
                    <span class="mt-2 text-sm font-medium" :class="[
                      currentStep >= index
                        ? 'text-brand-500'
                        : 'text-gray-500 dark:text-gray-400'
                    ]">
                      {{ step.name }}
                    </span>
                  </div>
                  <div
                    v-if="index < steps.length - 1"
                    class="hidden md:block w-24 h-1 mx-4"
                    :class="currentStep > index ? 'bg-brand-500' : 'bg-gray-200 dark:bg-slate-700'"
                  ></div>
                </div>
              </div>
            </div>
  
            <div class="grid lg:grid-cols-3 gap-8 px-3 py-3 md:px-8 md:py-8">

              <div class="lg:col-span-2">

                <div v-if="currentStep === 0" class="space-y-4">
                  <h2 class="text-2xl font-bold text-gray-900 dark:text-white mb-6">Review Your Order</h2>
                  <div class="space-y-4 w-screen md:w-full">
                    <CartItem
                      v-for="item in cart.items"
                      :key="item.menuItem.id"
                      :item="item"
                      @update-quantity="(change) => updateQuantity(item, change)"
                      @remove="removeItem(item)"
                    />
                  </div>
                  <div class="flex justify-between items-center pt-6 border-t">
                    <NuxtLink to="/menu" class="text-brand-500 hover:text-brand-600 font-medium flex items-center">
                      <ArrowLeft class="w-4 h-4 mr-2" />
                      Continue Shopping
                    </NuxtLink>
                    <button @click="currentStep = 1" class="btn-primary">
                      Continue to Delivery
                      <ArrowRight class="w-4 h-4 ml-2" />
                    </button>
                  </div>
                </div>
  

                <div v-if="currentStep === 1" class="space-y-8">
                  <OrderTypeSelector
                    :selected-type="orderData.orderType"
                    @select="orderData.orderType = $event"
                  />
                  
                  <CustomerInfoForm
                    @submit="handleCustomerInfo"
                    submit-text="Continue to Payment"
                  />
                  
                  <div class="flex justify-between pt-6 border-t">
                    <button @click="currentStep = 0" class="btn-secondary">
                      Back to Cart
                    </button>
                    <button @click="currentStep = 2" class="btn-primary">
                      Continue to Payment
                      <ArrowRight class="w-4 h-4 ml-2" />
                    </button>
                  </div>
                </div>
  

                <div v-if="currentStep === 2" class="space-y-8">
                  <div class="bg-white dark:bg-slate-800 rounded-xl p-6 border border-gray-200 dark:border-slate-700">
                    <h3 class="text-xl font-bold text-gray-900 dark:text-white mb-4">Payment Method</h3>
                    <div class="space-y-4">
                      <label
                        v-for="method in paymentMethods"
                        :key="method.id"
                        class="flex items-center p-4 border-2 rounded-lg cursor-pointer transition-all"
                        :class="[
                          selectedPayment === method.id
                            ? 'border-brand-500 bg-brand-50 dark:bg-brand-900/20'
                            : 'border-gray-200 dark:border-slate-700 hover:border-brand-300'
                        ]"
                      >
                        <input
                          v-model="selectedPayment"
                          type="radio"
                          :value="method.id"
                          class="mr-4"
                        />
                        <component :is="method.icon" class="w-6 h-6 mr-3" :class="method.iconClass" />
                        <div class="flex-1">
                          <h4 class="font-bold">{{ method.name }}</h4>
                          <p class="text-sm text-gray-600 dark:text-gray-400">{{ method.description }}</p>
                        </div>
                      </label>
                    </div>
                  </div>
  
                  <div class="flex justify-between pt-6 border-t">
                    <button @click="currentStep = 1" class="btn-secondary">
                      Back to Details
                    </button>
                    <button
                      @click="submitOrder"
                      :disabled="isSubmitting"
                      class="btn-primary"
                    >
                      <template v-if="isSubmitting">
                        <Loader2 class="w-5 h-5 animate-spin inline mr-2" />
                        Placing Order...
                      </template>
                      <template v-else>
                        Place Order - ${{ orderTotal.toFixed(2) }}
                      </template>
                    </button>
                  </div>
                </div>
              </div>
  

              <div class="lg:col-span-1">
                <CartSummary
                  :items="cart.items"
                  :delivery-fee="deliveryFee"
                  :is-loading="isSubmitting"
                  @checkout="currentStep = 1"
                />
              </div>
            </div>
          </div>
        </div>
      </div>
  

      <OrderConfirmation
        v-if="showConfirmation"
        :order-number="orderNumber"
        :estimated-time="estimatedTime"
        @close="showConfirmation = false"
      />
    </div>
  </template>
  
  <script setup lang="ts">
  import { ref, computed, onMounted } from 'vue'
  import { useCartStore } from '../../stores/cart'
  import type { Order, OrderType } from '../../types/order'
  import type { CartItem as CartItemType } from '../../types/menu'
  import {
    ShoppingCart,
    ShoppingBag,
    ArrowLeft,
    ArrowRight,
    Loader2,
    CreditCard,
    Smartphone,
    DollarSign
  } from 'lucide-vue-next'
  
  import PageHeader from '../../components/layout/PageHeader.vue'
  import CartItem from '../../components/order/CartItem.vue'
  import CartSummary from '../../components/order/CartSummary.vue'
  import OrderTypeSelector from '../../components/order/OrderTypeSelector.vue'
  import CustomerInfoForm from '../../components/order/CustomerInfoForm.vue'
  import OrderConfirmation from '../../components/order/OrderConfirmation.vue'
  
  const cart = useCartStore()
  
  const currentStep = ref(0)
  const isSubmitting = ref(false)
  const showConfirmation = ref(false)
  const orderNumber = ref('')
  const estimatedTime = ref('')
  const selectedPayment = ref('card')
  
  const orderData = ref<Partial<Order>>({
    orderType: 'pickup',
    customerName: '',
    customerPhone: '',
    customerEmail: '',
    notes: ''
  })
  
  const steps = [
    { id: 'cart', name: 'Cart' },
    { id: 'details', name: 'Details' },
    { id: 'payment', name: 'Payment' }
  ]
  
  const paymentMethods = [
    {
      id: 'card',
      name: 'Credit/Debit Card',
      description: 'Pay with Visa, MasterCard, or Amex',
      icon: CreditCard,
      iconClass: 'text-blue-500'
    },
    {
      id: 'mobile',
      name: 'Mobile Payment',
      description: 'Apple Pay, Google Pay',
      icon: Smartphone,
      iconClass: 'text-green-500'
    },
    {
      id: 'cash',
      name: 'Cash on Delivery',
      description: 'Pay when you receive your order',
      icon: DollarSign,
      iconClass: 'text-amber-500'
    }
  ]
  
  const deliveryFee = computed(() => {
    return orderData.value.orderType === 'delivery' ? 3.99 : 0
  })
  
  const orderTotal = computed(() => {
    const subtotal = cart.items.reduce((sum, item) => sum + (item.menuItem.price * item.quantity), 0)
    const tax = subtotal * 0.08
    return subtotal + tax + deliveryFee.value
  })
  
  const updateQuantity = (item: CartItemType, change: number) => {
    const newQuantity = item.quantity + change
    if (newQuantity > 0) {
      cart.updateQuantity(item.menuItem.id, newQuantity)
    } else {
      cart.removeItem(item.menuItem.id)
    }
  }
  
  const removeItem = (item: CartItemType) => {
    cart.removeItem(item.menuItem.id)
  }
  
  const handleCustomerInfo = (data: any) => {
    Object.assign(orderData.value, data)
  }
  
  const submitOrder = async () => {
    isSubmitting.value = true
  
    try {
      await new Promise(resolve => setTimeout(resolve, 1500))
  
      orderNumber.value = `FG-${Date.now().toString().slice(-6)}`
      
      const now = new Date()
      now.setMinutes(now.getMinutes() + 30)
      estimatedTime.value = now.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
  
      cart.clearCart()
  
      showConfirmation.value = true
      currentStep.value = 0
  
      orderData.value = {
        orderType: 'pickup',
        customerName: '',
        customerPhone: '',
        customerEmail: '',
        notes: ''
      }
  
    } catch (error) {
      console.error('Order submission failed:', error)
    } finally {
      isSubmitting.value = false
    }
  }
  
  onMounted(() => {
    window.scrollTo({ top: 0, behavior: 'smooth' })
  })
  </script>