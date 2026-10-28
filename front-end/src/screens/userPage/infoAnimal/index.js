import React from "react";
import { View, Text, TouchableOpacity, Image, ScrollView } from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";
import { Ionicons } from "@expo/vector-icons";

import styled from "./styled";

export default ({ navigation }) => {
  return (
    <SafeAreaView style={styled.safeArea}>
      {/* Cabeçalho */}
      <View style={styled.containerCabecalho}>
        <TouchableOpacity onPress={() => navigation?.goBack()}>
          <Ionicons name="arrow-back" size={28} color="#FFFFFF" />
        </TouchableOpacity>
        <Text style={styled.tituloCabecalho}>Informações do animal</Text>
      </View>

      {/* Cartão Principal com Rolagem (ScrollView) */}
      <ScrollView style={styled.suporteCartao} showsVerticalScrollIndicator={false}>
        <View style={styled.containerPrincipal}>
          
          {/* Bloco Superior */}
          <View style={styled.blocoSuperior}>
            <Image
              source={{
                uri: "https://www.infoescola.com/wp-content/uploads/2018/03/sagui-de-tufos-brancos_790108426-1000x667.jpg",
              }}
              style={styled.imagemAnimal}
            />

            <View style={styled.divisorVertical} />

            <View style={styled.containerInformacoes}>
              <Text style={styled.nomeAnimal}>SAGUI-DE-TUFOS-BRANCOS</Text>
              <View style={styled.divisorHorizontal} />
              <Text style={styled.nomeCientifico}>(Callithrix jacchus)</Text>
              {/* Linha do Peso Médio */}
              <View style={styled.linhaInformacao}>
                <Text style={styled.textoInformacao}>Peso médio:</Text>
                <Text style={styled.valorInformacao}>300 g - 450 g</Text>
              </View>
              {/* Linha do Tamanho Médio */}
              <View style={styled.linhaInformacao}>
                <Text style={styled.textoInformacao}>Tamanho médio:</Text>
                <Text style={styled.valorInformacao}>18 cm - 25 cm</Text>
              </View>
            </View>
          </View>

          {/* Bloco de Descrição */}
          <View style={styled.blocoDescricao}>
            <Text style={styled.tituloDescricao}>Descrição do animal</Text>
            <Text style={styled.textoDescricao}>
              O sagui-de-tufos-brancos é um primata de pequeno porte nativo do Nordeste do Brasil. Caracteriza-se pelos marcantes tufos brancos ao redor das orelhas e cauda anelada.
            </Text>
          </View>

        </View>
        
      </ScrollView>
      {/* Logo Discovery Rural no fim da pagina */}
        <Image
          style={{ width: 128, height: 128, alignSelf: 'center', marginTop: 20, marginBottom: 20 }}
          source={{ uri: "https://github.com/rafael9577/discoveryRural/blob/main/front-end/assets/logoDiscoveryRuralWhite.png?raw=true" }}
        />
    </SafeAreaView>
  );
};