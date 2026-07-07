// Базовая серверная логика для Zuhro.ai
const ALLOWED_DOMAINS = ["*"]; // Разрешить запросы со всех доменов

addEventListener('fetch', event => {
  event.respondWith(handleRequest(event.request))
})

async function handleRequest(request) {
  return new Response('Сервер Zuhro.ai запущен и работает!', {
    status: 200,
    headers: { 'content-type': 'text/plain' }
  })
}
