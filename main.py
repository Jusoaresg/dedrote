from twitchio.ext import commands
from obswebsocket import obsws,requests


PREFIX = "!"
TOKEN = "ylohuzwui7wl1czu14e0u8apj0ci7o"
CHANNEL_NAME = ""

SCENE = "Teste"
TEXT_SOURCE = "Texto"

OBS_HOST = ""
OBS_PORT = 0
OBS_PASSWORD = ""

class ObsWebSocket:
    ws = None

    def __init__(self):
        self.ws = obsws(OBS_HOST, OBS_PORT, OBS_PASSWORD)        

        try:
            self.ws.connect()
        except ConnectionError as e:
            print(f"PANIC: Failed to connect to OBS: {e}") 


        print("Conectou no OBS")

    def disconnect(self):
        if(self.ws is not None):
            self.ws.disconnect()

    def set_scene(self):
        if(self.ws is not None):
            self.ws.call(requests.SetCurrentProgramScene(sceneName=SCENE))

    def get_text(self, source_name):
        if(self.ws is not None):
            response = self.ws.call(requests.GetInputSettings(inputName=source_name))
            return response.datain["inputSettings"]["text"]

    async def set_text(self, source_name, new_text):
        if(self.ws is not None):
            self.ws.call(requests.SetInputSettings(inputName=source_name, inputSettings = {"text": new_text}))


class Bot(commands.Bot):
    obsSocket = None

    def __init__(self):
        super().__init__(
            token=TOKEN,
            prefix=PREFIX,
            initial_channels=[CHANNEL_NAME],
        )
        self.obsSocket = ObsWebSocket()

    async def event_ready(self):
        print(f'Logado como | {self.nick}')
        print(f'Id do usuario | {self.user_id}')

    @commands.command(aliases=("mortes", "morte"))
    async def death_count(self, ctx: commands.Context, message):

        if(self.obsSocket is not None and
           ctx.message.author.is_mod):

            await ctx.send("Contador mudado")
            await self.obsSocket.set_text(TEXT_SOURCE, f"Mortes: {message}")

bot = Bot()
bot.run()
