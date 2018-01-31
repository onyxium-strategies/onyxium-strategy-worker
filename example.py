from bittrex import Bittrex
import smtplib
import pprint
from dateutil.parser import parse
from datetime import datetime
import time

def email_sender(input_message, subject, email_to):
    ''' function to send email '''
    to = email_to
    gmail_user = 'alainfh94@gmail.com' ## email of sender account
    gmail_pwd = 'gfqodfvhlfwdmvft' ## password of sender account
    smtpserver = smtplib.SMTP("smtp.gmail.com",587)
    smtpserver.ehlo()
    smtpserver.starttls()
    smtpserver.ehlo
    smtpserver.login(gmail_user, gmail_pwd)
    header = 'To:' + to + '\n' + 'From: ' + gmail_user + '\n' +'Subject: ' + subject + '\n'
    msg = header + input_message
    smtpserver.sendmail(gmail_user, to, msg)
    smtpserver.close()

class altCoinTrader(object):
	def __init__(self, key, secret, altA, altB, threshold = 0, quantityToSell = 0):
		self.marketA = 'BTC-' + altA
		self.marketB = 'BTC-' + altB
		self.priceA = 0
		self.priceB = 0
		self.threshold = threshold
		self.quantityToSell = quantityToSell
		self.client = Bittrex(key, secret, calls_per_second=1)

	def trade(self):
		while True:
			if self.thresholdMet():
				msg = self.marketA + ' and ' + self.marketB + ' have a relative price of ' + str(self.priceA/self.priceB) + ' which is above the threshold of ' + str(self.threshold)
				email_sender(msg, 'Threshold met!', 'alainfh94@gmail.com')
				exit()

				# sellOrder = self.client.sell_limit(market=self.marketA, quantity=self.quantityToSell, rate=self.priceA)['result']
				# pp.pprint(sellOrder)
				# orderUuid = sellOrder['uuid']

				# self.waitForOrder(orderUuid)

				# # update priceB because it might have been a while since we last checked
				# currPriceB = self.client.get_ticker(self.marketB)['result']['Last']
				# if self.priceA/currPriceB > self.threshold:
				# 	self.priceB = currPriceB

				# btc = order['Price'] - order['CommissionPaid']
				# quantityToBuy = btc / self.priceB

				# buyOrder = self.client.buy_limit(market=self.marketB, quantity=quantityToBuy, rate=self.priceB)['result']
				# pp.pprint(buyOrder)
				# orderUuid = buyOrder['uuid']
				# finishedOrder = self.waitForOrder(orderUuid)

				# date = parse(finishedOrder['Closed'])
				# msg = 'Transaction has been completed on: ' + date.strftime('%d-%b-%Y %H:%M:%S')
				# email_sender(msg, 'alainfh94@gmail.com')
				# exit()

	def thresholdMet(self):
		self.priceA = self.client.get_ticker(self.marketA)['result']['Last']
		self.priceB = self.client.get_ticker(self.marketB)['result']['Last']
		pp.pprint(self.priceA)
		pp.pprint(self.priceB)
		pp.pprint(self.priceA/self.priceB)

		if self.priceA/self.priceB > self.threshold:
			# make sure the threshold is not matched by a spoof buy/sell
			# check if it is still valid after a minute.
			time.sleep(60)
			self.priceA = self.client.get_ticker(self.marketA)['result']['Last']
			self.priceB = self.client.get_ticker(self.marketB)['result']['Last']
			pp.pprint(self.priceA)
			pp.pprint(self.priceB)
			pp.pprint(self.priceA/self.priceB)
			if self.priceA/self.priceB > self.threshold:
				return True
		return False

	def waitForOrder(self, orderUuid):
		order = self.client.get_order(orderUuid)['result']
		pp.pprint(order)

		while order['IsOpen']:
			order = self.client.get_order(orderUuid)['result']
			pp.pprint(order)
		return order

def main():
	key = '48052a29012d4576ab92672eb938bac3'
	secret = 'f546b20541be4aea9cf9bb60b4ebc011'
	trader = altCoinTrader(key, secret, 'OMG', 'IOC', 2.0, 1)
	trader.trade()

if __name__ == '__main__':
	pp = pprint.PrettyPrinter(indent=4)
	main()
