var app = new Vue({
    el: '#app',
    data: {
        symbols: [
            "EURUSD",
            "USDJPY",
            "GBPUSD",
            "USDCHF",
            "EURCHF",
            "AUDUSD",
            "USDCAD",
            "NZDUSD",
            "EURGBP",
            "EURJPY",
            "GBPJPY",
            "CHFJPY",
        ],
        selectedSymbol: null,
        quotes:[]
    },
    methods: {
        getQuotes: function () {
            console.log('selectedSymbol', this.selectedSymbol)
            this.$http.post('http://localhost:9090/getQuotes', { symbol: this.selectedSymbol }).then(res => {
                console.log('res:', res)
                this.quotes.unshift(res.data.quotes[0])
            }, response => {
                console.log('error:', res)
            });
        }
    }
})