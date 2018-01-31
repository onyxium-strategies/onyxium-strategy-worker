from bittrex import Bittrex
import pprint # debug
import time
import argparse

def main(args):
	pass

if __name__ == '__main__':
	parser = argparse.ArgumentParser(description='Process a tree with trading conditions and actions.')
	parser.add_argument('input', nargs='?', help='json input file of tree', type=argparse.FileType('r'))
	parser.add_argument('exchange', help='exchange to connect with', choices=['bittrex', 'binance'] ,default='bittrex')
	parser.add_argument('key', nargs='?', help='api key', type=str)
	parser.add_argument('secret', nargs='?', help='api secret', type=str)
	args = vars(parser.parse_args())
	main(args)
