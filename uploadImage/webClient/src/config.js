export default {
  apiUrl: () => process.env.NODE_ENV === 'development' ? 'http://localhost:3083' : 'https://example.ru',
}
